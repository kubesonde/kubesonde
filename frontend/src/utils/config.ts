import { getRawApiServer } from "./env";

/**
 * Normalize the configured API server URL by stripping any trailing slash(es).
 * "http://x:2709/" -> "http://x:2709". Empty/unset values are returned as-is.
 */
const normalizeApiServer = (value: string | undefined): string | undefined =>
    value ? value.replace(/\/+$/, "") : value;

/**
 * Base URL of the Kubesonde controller REST API, taken from
 * `VITE_API_SERVER`. Undefined/empty when the app is not in API mode.
 */
export const apiServer: string | undefined = normalizeApiServer(getRawApiServer());

/**
 * True when a controller API server is configured, i.e. the UI should fetch
 * probe data live instead of relying on a manually uploaded file.
 */
export const isApiMode: boolean = !!apiServer;
