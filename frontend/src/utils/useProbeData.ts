import { useCallback, useEffect, useRef, useState } from "react";
import { apiServer } from "./config";
import { ProbeOutput } from "src/entities/probeOutput";

/** Interval between `/probes/status` polls, in milliseconds. */
export const POLL_INTERVAL_MS = 2_000;

/**
 * Shape of the `GET /probes/status` response.
 *
 * `window` is a Go `time.Duration` serialized as an integer number of
 * nanoseconds. We carry it through untouched (no interpretation here).
 */
export interface ProbeStatus {
    complete: boolean;
    count: number;
    window: number;
    secondsSinceLastChange: number;
}

export interface UseProbeDataResult {
    /** The full probe output, populated once the run is complete. */
    data: ProbeOutput | undefined;
    /** The latest status snapshot from `/probes/status`. */
    status: ProbeStatus | undefined;
    /** True once the run is complete and `data` has been loaded. */
    ready: boolean;
    /** True while a request is in flight. */
    loading: boolean;
    /** Set when a fetch fails or returns a non-2xx response. */
    error: string | undefined;
    /** Re-fetch immediately; restarts polling if not yet complete. */
    refresh: () => void;
}

/** GET a URL and parse JSON, throwing on a non-2xx response. */
const getJson = async <T>(url: string, signal: AbortSignal): Promise<T> => {
    // `no-store` so we always see the controller's current state rather than a
    // heuristically cached response.
    const resp = await fetch(url, { signal, cache: "no-store" });
    if (!resp.ok) {
        throw new Error(`request failed: ${resp.status}`);
    }
    return (await resp.json()) as T;
};

/** Fetch the current completeness status. */
export const fetchProbeStatus = (signal: AbortSignal): Promise<ProbeStatus> =>
    getJson<ProbeStatus>(`${apiServer}/probes/status`, signal);

/** Fetch the full probe output. */
export const fetchProbes = (signal: AbortSignal): Promise<ProbeOutput> =>
    getJson<ProbeOutput>(`${apiServer}/probes`, signal);

/**
 * A single poll snapshot: the current status and the full probe output.
 * Keeps all network logic out of the hook's state handling.
 */
export interface ProbeSnapshot {
    status: ProbeStatus;
    data: ProbeOutput;
}

/** Fetch the current status and the full probe output together. */
export const fetchProbeSnapshot = async (
    signal: AbortSignal
): Promise<ProbeSnapshot> => {
    const [status, data] = await Promise.all([
        fetchProbeStatus(signal),
        fetchProbes(signal),
    ]);
    return { status, data };
};

/**
 * Polls the Kubesonde controller for probe results.
 *
 * Every {@link POLL_INTERVAL_MS} it fetches a {@link fetchProbeSnapshot} (status
 * + full data). It publishes a new `data` reference only when the payload
 * actually changed, so the graph re-renders as new probes arrive but is left
 * untouched (preserving dragged node positions) when nothing changed. Polling
 * continues while mounted; `status.complete` is exposed as informational state.
 * Overlapping requests are suppressed and in-flight requests are aborted on
 * unmount.
 */
export const useProbeData = (): UseProbeDataResult => {
    const [data, setData] = useState<ProbeOutput | undefined>(undefined);
    const [status, setStatus] = useState<ProbeStatus | undefined>(undefined);
    const [ready, setReady] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | undefined>(undefined);

    // Interval handle for the poll timer.
    const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
    // AbortController for the in-flight request, so we can cancel on unmount.
    const controllerRef = useRef<AbortController | null>(null);
    // Guards against overlapping polls (don't start a new one while pending).
    const inFlightRef = useRef(false);
    // Tracks mounted state so async continuations don't touch state after unmount.
    const mountedRef = useRef(true);
    // True once we have loaded data at least once, so background refreshes don't
    // re-trigger the initial loading state.
    const hasDataRef = useRef(false);
    // Serialized last payload, so we only publish new `data` (a new reference)
    // when the probe results actually changed — otherwise the graph would
    // rebuild every poll and reset any nodes the user dragged.
    const lastDataJsonRef = useRef<string | undefined>(undefined);

    const clearTimer = useCallback(() => {
        if (intervalRef.current !== null) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
        }
    }, []);

    // Performs one poll cycle, delegating all network work to fetchProbeSnapshot.
    const poll = useCallback(async () => {
        if (inFlightRef.current) {
            return;
        }
        inFlightRef.current = true;

        const controller = new AbortController();
        controllerRef.current = controller;

        // Only surface the loading state before we have any data; background
        // refreshes update the graph silently.
        if (mountedRef.current && !hasDataRef.current) {
            setLoading(true);
        }

        try {
            const { status: nextStatus, data: nextData } = await fetchProbeSnapshot(
                controller.signal
            );
            if (!mountedRef.current) {
                return;
            }
            hasDataRef.current = true;
            setStatus(nextStatus);
            // Only publish a new data reference when the payload changed, so the
            // graph does not rebuild (and lose dragged positions) every poll.
            const nextJson = JSON.stringify(nextData);
            if (nextJson !== lastDataJsonRef.current) {
                lastDataJsonRef.current = nextJson;
                setData(nextData);
            }
            setReady(true);
            setError(undefined);
        } catch (err) {
            // Ignore aborts caused by unmount / refresh; surface everything else.
            if (controller.signal.aborted || !mountedRef.current) {
                return;
            }
            setError(
                err instanceof Error ? err.message : "Failed to fetch probe data"
            );
        } finally {
            inFlightRef.current = false;
            if (mountedRef.current) {
                setLoading(false);
            }
        }
    }, []);

    const startPolling = useCallback(() => {
        clearTimer();
        intervalRef.current = setInterval(() => {
            void poll();
        }, POLL_INTERVAL_MS);
    }, [clearTimer, poll]);

    const refresh = useCallback(() => {
        // Cancel any in-flight request so the immediate fetch is authoritative,
        // then poll now. The background interval keeps running regardless.
        if (controllerRef.current) {
            controllerRef.current.abort();
        }
        inFlightRef.current = false;
        void poll();
    }, [poll]);

    useEffect(() => {
        mountedRef.current = true;
        // Kick off an immediate poll and start the interval.
        void poll();
        startPolling();

        return () => {
            mountedRef.current = false;
            clearTimer();
            if (controllerRef.current) {
                controllerRef.current.abort();
            }
            // Release the guard so a remount (e.g. React StrictMode) can poll
            // immediately instead of waiting for the aborted request to settle.
            inFlightRef.current = false;
        };
        // Poll/startPolling/clearTimer are stable (memoized); run once on mount.
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return { data, status, ready, loading, error, refresh };
};

export default useProbeData;
