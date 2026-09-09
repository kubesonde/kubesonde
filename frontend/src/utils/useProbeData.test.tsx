import { act, renderHook, waitFor } from "@testing-library/react";
import { POLL_INTERVAL_MS, useProbeData } from "./useProbeData";

// `config.ts` reads env at load time; pin apiServer for URL construction.
jest.mock("./config", () => ({
    apiServer: "http://ctrl:2709",
    isApiMode: true,
}));

type FetchMock = jest.MockedFunction<typeof fetch>;

const jsonResponse = (body: unknown, ok = true, statusCode = 200): Response =>
    ({
        ok,
        status: statusCode,
        json: async () => body,
    } as unknown as Response);

const incompleteStatus = { complete: false, count: 3, window: 1000, secondsSinceLastChange: 1 };
const completeStatus = { complete: true, count: 5, window: 1000, secondsSinceLastChange: 30 };
const probeOutput = { start: "a", end: "b", items: [], errors: [], podNetworkingv2: {}, podConfigurationNetworking: {} };

describe("useProbeData", () => {
    let fetchMock: FetchMock;

    beforeEach(() => {
        jest.useFakeTimers();
        fetchMock = jest.fn() as unknown as FetchMock;
        globalThis.fetch = fetchMock;
    });

    afterEach(() => {
        jest.runOnlyPendingTimers();
        jest.useRealTimers();
        jest.clearAllMocks();
    });

    // Respond based on URL: status endpoint vs full probes payload.
    const mockByUrl = (status: unknown, probes: unknown) => {
        fetchMock.mockImplementation((input: RequestInfo | URL) => {
            const url = String(input);
            return Promise.resolve(
                url.endsWith("/probes/status")
                    ? jsonResponse(status)
                    : jsonResponse(probes)
            );
        });
    };

    it("renders data immediately without waiting for complete, and keeps polling", async () => {
        mockByUrl(incompleteStatus, probeOutput);

        const { result } = renderHook(() => useProbeData());

        // First poll fetches BOTH status and /probes, exposing data even though
        // the run is not complete.
        await waitFor(() => expect(result.current.ready).toBe(true));
        expect(result.current.status?.complete).toBe(false);
        expect(result.current.data).toEqual(probeOutput);
        expect(fetchMock).toHaveBeenCalledWith(
            "http://ctrl:2709/probes/status",
            expect.anything()
        );
        expect(fetchMock).toHaveBeenCalledWith(
            "http://ctrl:2709/probes",
            expect.anything()
        );

        // Polling continues (never stops), so more fetches happen over time.
        const before = fetchMock.mock.calls.length;
        await act(async () => {
            jest.advanceTimersByTime(POLL_INTERVAL_MS);
        });
        await waitFor(() =>
            expect(fetchMock.mock.calls.length).toBeGreaterThan(before)
        );
    });

    it("keeps polling even after the run is complete (late probes still arrive)", async () => {
        mockByUrl(completeStatus, probeOutput);

        const { result } = renderHook(() => useProbeData());

        await waitFor(() => expect(result.current.ready).toBe(true));
        expect(result.current.status?.complete).toBe(true);

        const before = fetchMock.mock.calls.length;

        // Polling continues so edges recorded after completion are still fetched.
        await act(async () => {
            jest.advanceTimersByTime(POLL_INTERVAL_MS);
        });
        await waitFor(() =>
            expect(fetchMock.mock.calls.length).toBeGreaterThan(before)
        );
    });

    it("does not publish a new data reference when the payload is unchanged", async () => {
        mockByUrl(incompleteStatus, probeOutput);

        const { result } = renderHook(() => useProbeData());
        await waitFor(() => expect(result.current.data).toBeDefined());
        const firstData = result.current.data;

        // Another poll with identical payload must keep the same reference so the
        // graph is not rebuilt (dragged node positions are preserved).
        await act(async () => {
            jest.advanceTimersByTime(POLL_INTERVAL_MS);
        });
        await waitFor(() =>
            expect(fetchMock.mock.calls.length).toBeGreaterThan(2)
        );
        expect(result.current.data).toBe(firstData);
    });

    it("refresh() re-fetches immediately", async () => {
        fetchMock.mockResolvedValue(jsonResponse(incompleteStatus));

        const { result } = renderHook(() => useProbeData());
        await waitFor(() => expect(result.current.status).toBeDefined());

        const before = fetchMock.mock.calls.length;
        await act(async () => {
            result.current.refresh();
        });
        await waitFor(() =>
            expect(fetchMock.mock.calls.length).toBeGreaterThan(before)
        );
    });

    it("sets error state on a failed fetch", async () => {
        fetchMock.mockResolvedValue(jsonResponse(null, false, 503));

        const { result } = renderHook(() => useProbeData());

        await waitFor(() => expect(result.current.error).toBeDefined());
        expect(result.current.ready).toBe(false);
    });

    it("clears interval and aborts in-flight fetches on unmount", async () => {
        const abortSpy = jest.spyOn(AbortController.prototype, "abort");
        fetchMock.mockResolvedValue(jsonResponse(incompleteStatus));

        const { result, unmount } = renderHook(() => useProbeData());
        await waitFor(() => expect(result.current.status).toBeDefined());

        const before = fetchMock.mock.calls.length;
        unmount();

        expect(abortSpy).toHaveBeenCalled();

        // No further polls after unmount.
        await act(async () => {
            jest.advanceTimersByTime(POLL_INTERVAL_MS * 2);
        });
        expect(fetchMock.mock.calls.length).toBe(before);

        abortSpy.mockRestore();
    });
});
