import React from "react";
import { shallow } from "enzyme";
import Row from ".";

describe(Row, () => {
  it("renders the content correctly", () => {
    const text = "My Text";
    const wrapper = shallow(<Row>{text}</Row>);

    expect(wrapper).toHaveText(text);
    expect(wrapper).toMatchSnapshot();
  });

  it("includes the expected CSS class", () => {
    const wrapper = shallow(<Row>Test</Row>);
    expect(wrapper).toHaveClassName("clr-row");
  });

  it("add the role when the row is a list", () => {
    const wrapper = shallow(<Row list>Test</Row>);
    expect(wrapper.prop("role")).toBe("list");
  });
});
