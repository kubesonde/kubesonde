import { useCallback, useEffect, useRef, useState } from "react";
import { apiServer } from "./config";
import { ProbeOutput } from "src/entities/probeOutput";

/** Interval between `/probes/status` polls, in milliseconds. */
export const POLL_INTERVAL_MS = 10_000;

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

const statusUrl = () => `${apiServer}/probes/status`;
const probesUrl = () => `${apiServer}/probes`;

/**
 * Polls the Kubesonde controller for probe results.
 *
 * Every {@link POLL_INTERVAL_MS} it fetches `/probes/status`. While the run is
 * incomplete it keeps polling and exposes `status.count` for progress display.
 * Once `status.complete` is true it fetches the full `/probes` payload, stops
 * the interval (frontend-only; the controller keeps probing), and exposes the
 * data. Overlapping requests are suppressed and in-flight requests are aborted
 * on unmount.
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

    const clearTimer = useCallback(() => {
        if (intervalRef.current !== null) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
        }
    }, []);

    // Performs one poll cycle: status, then /probes when complete.
    const poll = useCallback(async () => {
        if (inFlightRef.current) {
            return;
        }
        inFlightRef.current = true;

        const controller = new AbortController();
        controllerRef.current = controller;

        if (mountedRef.current) {
            setLoading(true);
        }

        try {
            const statusResp = await fetch(statusUrl(), { signal: controller.signal });
            if (!statusResp.ok) {
                throw new Error(`status request failed: ${statusResp.status}`);
            }
            const nextStatus = (await statusResp.json()) as ProbeStatus;
            if (!mountedRef.current) {
                return;
            }
            setStatus(nextStatus);

            if (nextStatus.complete) {
                const probesResp = await fetch(probesUrl(), { signal: controller.signal });
                if (!probesResp.ok) {
                    throw new Error(`probes request failed: ${probesResp.status}`);
                }
                const nextData = (await probesResp.json()) as ProbeOutput;
                if (!mountedRef.current) {
                    return;
                }
                setData(nextData);
                setReady(true);
                setError(undefined);
                // Run is complete: stop polling (frontend-only).
                clearTimer();
            } else {
                if (!mountedRef.current) {
                    return;
                }
                setReady(false);
                setError(undefined);
            }
        } catch (err) {
            // Ignore aborts caused by unmount / refresh; surface everything else.
            if (controller.signal.aborted) {
                return;
            }
            if (mountedRef.current) {
                setError(
                    err instanceof Error ? err.message : "Failed to fetch probe data"
                );
            }
        } finally {
            inFlightRef.current = false;
            if (mountedRef.current) {
                setLoading(false);
            }
        }
    }, [clearTimer]);

    const startPolling = useCallback(() => {
        clearTimer();
        intervalRef.current = setInterval(() => {
            void poll();
        }, POLL_INTERVAL_MS);
    }, [clearTimer, poll]);

    const refresh = useCallback(() => {
        // Cancel any in-flight request so the immediate fetch is authoritative.
        if (controllerRef.current) {
            controllerRef.current.abort();
        }
        inFlightRef.current = false;
        // Restart polling if we haven't completed yet.
        if (!ready) {
            startPolling();
        }
        void poll();
    }, [poll, ready, startPolling]);

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
        };
        // Poll/startPolling/clearTimer are stable (memoized); run once on mount.
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return { data, status, ready, loading, error, refresh };
};

export default useProbeData;
