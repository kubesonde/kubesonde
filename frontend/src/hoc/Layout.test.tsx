import React from "react";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";

// isApiMode is controlled per-test by mutating this mock's property.
jest.mock("src/utils/config", () => ({
  __esModule: true,
  isApiMode: false,
  apiServer: undefined,
}));

// Heavy children are stubbed so the test only exercises routing.
jest.mock("src/components/graph/GraphFromApi", () => ({
  __esModule: true,
  GraphFromApi: () => <div>LIVE VIEW STUB</div>,
  default: () => <div>LIVE VIEW STUB</div>,
}));
jest.mock("src/components/graph/GraphJSONUploadComponent", () => ({
  __esModule: true,
  GraphJSONUploadComponent: () => <div>UPLOAD STUB</div>,
}));
jest.mock("src/components/graph/ExampleGraph", () => ({
  __esModule: true,
  ExampleGraphComponent: () => <div>EXAMPLE STUB</div>,
}));
jest.mock("src/components/graph/GraphFromLocation", () => ({
  __esModule: true,
  GraphFromLocation: () => <div>GRAPH FROM LOCATION STUB</div>,
}));
// Sidebar relies on import.meta / react-pro-sidebar and is not under test here.
jest.mock("./Sidebar", () => ({
  __esModule: true,
  Sidebar: () => <div>SIDEBAR STUB</div>,
}));

import { Layout } from "./Layout";
// eslint-disable-next-line @typescript-eslint/no-var-requires
const config = require("src/utils/config") as { isApiMode: boolean };

const renderAt = (path: string) =>
  render(
    <MemoryRouter initialEntries={[path]}>
      <Layout />
    </MemoryRouter>
  );

describe("Layout landing route", () => {
  it("renders the live view at / when isApiMode is true", () => {
    config.isApiMode = true;
    renderAt("/");
    expect(screen.getByText("LIVE VIEW STUB")).toBeInTheDocument();
  });

  it("renders HomeComponent at / when isApiMode is false", () => {
    config.isApiMode = false;
    renderAt("/");
    // HomeComponent renders the Kubesonde hero title.
    expect(screen.getByText("Kubesonde")).toBeInTheDocument();
    expect(screen.queryByText("LIVE VIEW STUB")).not.toBeInTheDocument();
  });
});
