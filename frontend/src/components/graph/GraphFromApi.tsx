import React, { useEffect, useState } from "react";
import { FiRefreshCw } from "react-icons/fi";
import { Button } from "src/components/button/Button";
import { GraphView } from "src/components/graph/GraphView";
import { useProbeData } from "src/utils/useProbeData";
import "./GraphFromApi.scss";

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
  const { data, status, ready, loading, error, refresh } = useProbeData();

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
        <div className="live-toolbar">
          <span className="live-toolbar__status">
            <span
              className={
                "live-toolbar__dot" +
                (status?.complete ? "" : " live-toolbar__dot--live")
              }
            />
            {status?.complete ? "Probing complete" : "Probing…"}
            {status ? ` · ${status.count} probes` : ""}
            {lastUpdated
              ? ` · updated ${lastUpdated.toLocaleTimeString()}`
              : ""}
          </span>
          <Button
            title={loading ? "Refreshing…" : "Refresh"}
            icon={
              <FiRefreshCw
                className={loading ? "live-toolbar__spin" : undefined}
              />
            }
            onClick={() => refresh()}
          />
        </div>
        <GraphView data={data} title="Live cluster" />
      </>
    );
  }

  if (error) {
    return (
      <div className="live-state">Probes are not ready or not available</div>
    );
  }

  return (
    <div className="live-state">
      <span className="live-state__spinner" />
      Probing… {status?.count ?? 0} probes so far
    </div>
  );
};

export default GraphFromApi;
