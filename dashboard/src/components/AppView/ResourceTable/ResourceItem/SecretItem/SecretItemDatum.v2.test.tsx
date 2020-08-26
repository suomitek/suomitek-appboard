import { mount } from "enzyme";
import * as React from "react";

import { CdsIcon } from "components/Clarity/clarity";
import SecretItemDatum from "./SecretItemDatum.v2";

const testProps = {
  name: "foo",
  value: "YmFy", // foo
};

it("renders the secret datum (hidden by default)", () => {
  const wrapper = mount(<SecretItemDatum {...testProps} />);
  expect(wrapper.find(CdsIcon).findWhere(i => i.prop("shape") === "eye")).toExist();
  expect(wrapper.find(CdsIcon).findWhere(i => i.prop("shape") === "copy-to-clipboard")).toExist();
  expect(wrapper).toMatchSnapshot();
});

it("displays the secret datum value when clicking on the icon", () => {
  const wrapper = mount(<SecretItemDatum {...testProps} />);
  expect(wrapper.find("input").props()).toMatchObject({
    type: "password",
    value: "bar",
  });
  const icon = wrapper.find("button").findWhere(b => b.prop("aria-label") === "Show Secret");
  icon.simulate("click");
  wrapper.update();
  expect(wrapper.find(CdsIcon).findWhere(i => i.prop("shape") === "eye-hide")).toExist();
  expect(wrapper.find("input").props()).toMatchObject({
    type: "text",
    value: "bar",
  });
});
