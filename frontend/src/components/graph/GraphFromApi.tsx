import React, { useEffect, useState } from "react";
import { GraphView } from "src/components/graph/GraphView";
import { useProbeData } from "src/utils/useProbeData";

/**
 * Container that drives {@link GraphView} from the live controller API.
 *
 * It fetches its own data via {@link useProbeData} (state, not props) and:
 * - while loading / not yet ready (and no error): shows a progress indicator;
 * - on error, or when the controller is unreachable / not ready with no data:
 *   shows a simple "not available" message and lets the hook keep retrying;
 * - once ready with data: renders `<GraphView data title="Live cluster" />`.
 *
 * The manual refresh button (Task 6) is intentionally not wired here; it will be
 * slotted in alongside the ready view using the hook's `refresh()`.
 */
export const GraphFromApi = () => {
  const { data, status, ready, error, refresh } = useProbeData();

  // Track when the probe data was last loaded so we can show a "last updated"
  // timestamp next to the refresh button. Updated whenever `data` changes.
  const [lastUpdated, setLastUpdated] = useState<Date | undefined>(undefined);

  useEffect(() => {
    if (data) {
      setLastUpdated(new Date());
    }
  }, [data]);

  if (ready && data) {
    return (
      <>
        {/* Task 6: manual refresh button, using the hook's refresh(). */}
        <div className="refresh-bar">
          <button
            type="button"
            className="btn btn-primary"
            onClick={() => refresh()}
          >
            Refresh
          </button>
          {lastUpdated && (
            <span className="last-updated">
              Last updated: {lastUpdated.toLocaleTimeString()}
            </span>
          )}
        </div>
        <GraphView data={data} title="Live cluster" />
      </>
    );
  }

  if (error) {
    return <div>Probes are not ready or not available</div>;
  }

  return <div>Probing… {status?.count ?? 0} probes so far</div>;
};

export default GraphFromApi;
