import { render, screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom";
import React from "react";
import { GraphFromApi } from "./GraphFromApi";
import { useProbeData, UseProbeDataResult } from "src/utils/useProbeData";
import { ProbeOutput } from "src/entities/probeOutput";

// useProbeData is the container's only data source. Mock it with a factory so
// each test can drive a specific hook state (loading / error / ready). A factory
// avoids loading the real module (which chains to `import.meta.env` via config).
jest.mock("src/utils/useProbeData", () => ({
  useProbeData: jest.fn(),
}));

// GraphView pulls in GraphBase/Cytoscape which is hard to render in jsdom.
// Stub it and surface the title so we can assert the container rendered it
// with the data it received.
jest.mock("src/components/graph/GraphView", () => ({
  GraphView: (props: { data: unknown; title: string }) => (
    <div data-testid="graph-view">{props.title}</div>
  ),
}));

const mockedUseProbeData = useProbeData as jest.MockedFunction<
  typeof useProbeData
>;

const baseResult: UseProbeDataResult = {
  data: undefined,
  status: undefined,
  ready: false,
  loading: false,
  error: undefined,
  refresh: jest.fn(),
};

const mockHook = (overrides: Partial<UseProbeDataResult>) => {
  mockedUseProbeData.mockReturnValue({ ...baseResult, ...overrides });
};

afterEach(() => {
  jest.clearAllMocks();
});

test("shows progress text while loading / not ready", () => {
  mockHook({
    loading: true,
    status: { complete: false, count: 7, window: 0, secondsSinceLastChange: 0 },
  });

  render(<GraphFromApi />);

  expect(screen.getByText(/Probing/)).toBeInTheDocument();
  expect(screen.getByText(/7 probes so far/)).toBeInTheDocument();
  expect(screen.queryByTestId("graph-view")).not.toBeInTheDocument();
});

test("defaults to 0 probes when no status is available yet", () => {
  mockHook({});

  render(<GraphFromApi />);

  expect(screen.getByText(/0 probes so far/)).toBeInTheDocument();
});

test("shows the not-available message on error", () => {
  mockHook({ error: "controller unreachable" });

  render(<GraphFromApi />);

  expect(
    screen.getByText("Probes are not ready or not available")
  ).toBeInTheDocument();
  expect(screen.queryByTestId("graph-view")).not.toBeInTheDocument();
});

test("renders GraphView with the Live cluster title once ready", () => {
  mockHook({
    ready: true,
    data: {} as ProbeOutput,
    status: { complete: true, count: 12, window: 0, secondsSinceLastChange: 30 },
  });

  render(<GraphFromApi />);

  expect(screen.getByTestId("graph-view")).toBeInTheDocument();
  expect(screen.getByText("Live cluster")).toBeInTheDocument();
});

test("does not show the refresh button until ready", () => {
  mockHook({
    loading: true,
    status: { complete: false, count: 3, window: 0, secondsSinceLastChange: 0 },
  });

  render(<GraphFromApi />);

  expect(
    screen.queryByText(/refresh/i)
  ).not.toBeInTheDocument();
  expect(screen.queryByText(/Last updated/i)).not.toBeInTheDocument();
});

test("does not show the refresh button on error", () => {
  mockHook({ error: "controller unreachable" });

  render(<GraphFromApi />);

  expect(
    screen.queryByText(/refresh/i)
  ).not.toBeInTheDocument();
});

test("shows the refresh button and last-updated text once ready", () => {
  mockHook({
    ready: true,
    data: {} as ProbeOutput,
    status: { complete: true, count: 12, window: 0, secondsSinceLastChange: 30 },
  });

  render(<GraphFromApi />);

  expect(
    screen.getByText(/refresh/i)
  ).toBeInTheDocument();
  expect(screen.getByText(/updated/i)).toBeInTheDocument();
});

test("clicking refresh calls the hook's refresh function", () => {
  const refresh = jest.fn();
  mockHook({
    ready: true,
    data: {} as ProbeOutput,
    status: { complete: true, count: 12, window: 0, secondsSinceLastChange: 30 },
    refresh,
  });

  render(<GraphFromApi />);

  fireEvent.click(screen.getByText(/refresh/i));

  expect(refresh).toHaveBeenCalledTimes(1);
});
