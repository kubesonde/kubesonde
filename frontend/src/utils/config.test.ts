// `tsconfig` limits `types` to vite/client, so `require` isn't globally typed
// here; declare it locally for the module-reset re-import pattern below.
declare const require: (module: string) => unknown;

// `config.ts` reads the env value at module load, so each case re-imports the
// module with `./env` mocked to a specific raw value.
const loadConfig = (raw: string | undefined) => {
    jest.resetModules();
    jest.doMock("./env", () => ({
        getRawApiServer: () => raw,
    }));
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    return require("./config") as typeof import("./config");
};

describe("config", () => {
    afterEach(() => {
        jest.resetModules();
        jest.dontMock("./env");
    });

    it("apiServer unset -> isApiMode false", () => {
        const { apiServer, isApiMode } = loadConfig(undefined);
        expect(apiServer).toBeFalsy();
        expect(isApiMode).toBe(false);
    });

    it("empty string -> isApiMode false", () => {
        const { apiServer, isApiMode } = loadConfig("");
        expect(apiServer).toBeFalsy();
        expect(isApiMode).toBe(false);
    });

    it("apiServer set -> isApiMode true", () => {
        const { apiServer, isApiMode } = loadConfig("http://x:2709");
        expect(apiServer).toBe("http://x:2709");
        expect(isApiMode).toBe(true);
    });

    it("strips a trailing slash", () => {
        const { apiServer } = loadConfig("http://x:2709/");
        expect(apiServer).toBe("http://x:2709");
    });

    it("strips multiple trailing slashes", () => {
        const { apiServer } = loadConfig("http://x:2709///");
        expect(apiServer).toBe("http://x:2709");
    });
});
