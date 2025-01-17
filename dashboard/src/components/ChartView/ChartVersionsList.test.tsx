import { shallow } from "enzyme";
import * as React from "react";

import { Link } from "react-router-dom";
import { IChartVersion } from "../../shared/types";
import * as url from "../../shared/url";
import ChartVersionsList from "./ChartVersionsList";

const testChart: IChartVersion["relationships"]["chart"] = {
  data: {
    name: "test",
    repo: {
      name: "testrepo",
      namespace: "suomitek-appboard",
    },
  },
} as IChartVersion["relationships"]["chart"];

const testVersions: IChartVersion[] = [
  {
    attributes: {
      created: "2016-10-19T00:03:14.037Z",
      version: "1.2.3",
    },
    id: "1",
    relationships: { chart: testChart },
  },
  {
    attributes: {
      created: "2016-10-19T00:03:14.037Z",
      version: "1.2.2",
    },
    id: "2",
    relationships: { chart: testChart },
  },
  {
    attributes: {
      created: "2016-10-19T00:03:14.037Z",
      version: "1.2.1",
    },
    id: "3",
    relationships: { chart: testChart },
  },
  {
    attributes: {
      created: "2016-10-19T00:03:14.037Z",
      version: "1.2.0",
    },
    id: "4",
    relationships: { chart: testChart },
  },
] as IChartVersion[];

const extendedVersions: IChartVersion[] = [
  ...testVersions,
  {
    attributes: { created: "2016-10-19T00:03:14.037Z", version: "1.1.9" },
    id: "5",
    relationships: { chart: testChart },
  },
  {
    attributes: { created: "2016-10-19T00:03:14.037Z", version: "1.1.8" },
    id: "6",
    relationships: { chart: testChart },
  },
  {
    attributes: { created: "2016-10-19T00:03:14.037Z", version: "1.1.7" },
    id: "7",
    relationships: { chart: testChart },
  },
] as IChartVersion[];

it("renders the list of versions", () => {
  const wrapper = shallow(
    <ChartVersionsList
      versions={testVersions}
      selected={testVersions[1]}
      targetNamespace="targetNamespace"
      cluster="default"
    />,
  );
  const items = wrapper.find("li");
  expect(items).toHaveLength(4);
  const link = items.at(1).find(Link);
  expect(link.prop("className")).toBe("type-bold type-color-action");
  // The link include the namespace
  expect(link.prop("to")).toBe(
    url.app.charts.version("default", "targetNamespace", "test", "1.2.2", testChart.data.repo),
  );
});

it("does not render a the Show All link when there are 5 or less versions", () => {
  const wrapper = shallow(
    <ChartVersionsList
      versions={testVersions}
      selected={testVersions[1]}
      targetNamespace="targetNamespace"
      cluster="default"
    />,
  );
  expect(wrapper.find("a").exists()).toBe(false);
  wrapper.setProps({
    versions: [
      {
        attributes: { created: "2016-10-19T00:03:14.037Z", version: "1.2.4" },
        id: "0",
        relationships: { chart: testChart },
      },
      ...testVersions,
    ],
  });
  expect(wrapper.find("a").exists()).toBe(false);
});

it("renders a the Show All link when there are more than 5 versions", () => {
  const wrapper = shallow(
    <ChartVersionsList
      versions={extendedVersions}
      selected={extendedVersions[1]}
      targetNamespace="targetNamespace"
      cluster="default"
    />,
  );
  const showAllLink = wrapper.find("button");
  expect(showAllLink.exists()).toBe(true);
  expect(showAllLink.text()).toBe("Show all...");
  const items = wrapper.find("li");
  expect(items).toHaveLength(5);
});

it("shows all the versions when the Show All link is clicked", () => {
  const wrapper = shallow(
    <ChartVersionsList
      versions={extendedVersions}
      selected={extendedVersions[1]}
      targetNamespace="targetNamespace"
      cluster="default"
    />,
  );
  expect(wrapper.find("li")).toHaveLength(5);
  wrapper.find("button").simulate("click");
  expect(wrapper.find("li")).toHaveLength(7);
});
