import * as React from "react";
import ReactMarkdown from "react-markdown";
import * as ReactRedux from "react-redux";
import { HashLink as Link } from "react-router-hash-link";

import actions from "actions";
import LoadingWrapper from "components/LoadingWrapper/LoadingWrapper.v2";
import { defaultStore, mountWrapper } from "shared/specs/mountWrapper";
import ChartReadme from "./ChartReadme.v2";

const namespace = "chart-namespace";
const version = "1.2.3";
const defaultProps = {
  namespace,
  version,
  chartID: "stable/wordpress",
  getChartReadme: jest.fn(),
};

let spyOnUseDispatch: jest.SpyInstance;
const kubeaActions = { ...actions.kube };
beforeEach(() => {
  actions.charts = {
    ...actions.charts,
    getChartReadme: jest.fn(),
  };
  const mockDispatch = jest.fn();
  spyOnUseDispatch = jest.spyOn(ReactRedux, "useDispatch").mockReturnValue(mockDispatch);
});

afterEach(() => {
  actions.kube = { ...kubeaActions };
  spyOnUseDispatch.mockRestore();
});

it("behaves as a loading component", () => {
  const wrapper = mountWrapper(defaultStore, <ChartReadme {...defaultProps} />);
  expect(wrapper.find(LoadingWrapper)).toExist();
});

describe("getChartReadme", () => {
  it("gets triggered when mounting", () => {
    mountWrapper(defaultStore, <ChartReadme {...defaultProps} />);
    expect(actions.charts.getChartReadme).toHaveBeenCalledWith(
      namespace,
      defaultProps.chartID,
      version,
    );
  });
});

it("renders the ReactMarkdown content is readme is present", () => {
  const props = {
    ...defaultProps,
    readme: "# Markdown Readme",
  };
  const wrapper = mountWrapper(defaultStore, <ChartReadme {...props} />);
  const component = wrapper.find(ReactMarkdown);
  expect(component.props()).toMatchObject({ source: "# Markdown Readme" });
});

it("renders a not found error when error is set", () => {
  const wrapper = mountWrapper(defaultStore, <ChartReadme {...defaultProps} error={"not found"} />);
  expect(wrapper.text()).toContain("No README found");
});

it("renders an alert when error is set", () => {
  const wrapper = mountWrapper(defaultStore, <ChartReadme {...defaultProps} error={"Boom!"} />);
  expect(wrapper.text()).toContain("Unable to fetch chart README: Boom!");
});

it("renders the ReactMarkdown content adding IDs for the titles", () => {
  const wrapper = mountWrapper(
    defaultStore,
    <ChartReadme {...defaultProps} readme="# _Markdown_ 'Readme_or_not'!" />,
  );
  const component = wrapper.find("#markdown-readme_or_not");
  expect(component).toExist();
});

it("renders the ReactMarkdown ignoring comments", () => {
  const wrapper = mountWrapper(
    defaultStore,
    <ChartReadme
      {...defaultProps}
      readme={`<!-- This is a comment -->
      This is text`}
    />,
  );
  const html = wrapper.html();
  expect(html).toContain("This is text");
  expect(html).not.toContain("This is a comment");
});

it("renders the ReactMarkdown content with hash links", () => {
  const wrapper = mountWrapper(
    defaultStore,
    <ChartReadme
      {...defaultProps}
      readme={`[section 1](#section-1)
      # Section 1`}
    />,
  );
  expect(wrapper.find(Link)).toExist();
});
