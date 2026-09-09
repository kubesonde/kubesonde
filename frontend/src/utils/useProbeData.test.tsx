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

    it("polls status until complete, then fetches /probes and stops the interval", async () => {
        // First poll: incomplete. Second poll: complete -> also fetches /probes.
        fetchMock
            .mockResolvedValueOnce(jsonResponse(incompleteStatus)) // status #1
            .mockResolvedValueOnce(jsonResponse(completeStatus)) // status #2
            .mockResolvedValueOnce(jsonResponse(probeOutput)); // /probes

        const { result } = renderHook(() => useProbeData());

        // Immediate poll on mount hits /probes/status.
        await waitFor(() => expect(result.current.status?.complete).toBe(false));
        expect(fetchMock).toHaveBeenCalledWith(
            "http://ctrl:2709/probes/status",
            expect.anything()
        );
        expect(result.current.ready).toBe(false);

        // Advance one interval: second status poll returns complete -> /probes.
        await act(async () => {
            jest.advanceTimersByTime(POLL_INTERVAL_MS);
        });

        await waitFor(() => expect(result.current.ready).toBe(true));
        expect(result.current.data).toEqual(probeOutput);
        expect(fetchMock).toHaveBeenCalledWith(
            "http://ctrl:2709/probes",
            expect.anything()
        );

        const callsAfterComplete = fetchMock.mock.calls.length;

        // Interval must be stopped: further time passes without new fetches.
        await act(async () => {
            jest.advanceTimersByTime(POLL_INTERVAL_MS * 3);
        });
        expect(fetchMock.mock.calls.length).toBe(callsAfterComplete);
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
