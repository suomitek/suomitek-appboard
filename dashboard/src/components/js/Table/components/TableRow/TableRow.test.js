import React from "react";
import { shallow } from "enzyme";
import TableRow from ".";

const columns = [
  {
    accessor: "uuid",
    Header: "UUID",
  },
  {
    accessor: "name",
    Header: "Name",
  },
];

const row = {
  uuid: "1khj1kjas-quhkjwa-qjkwdka-1dkjasdna",
  name: "test",
};

describe(TableRow, () => {
  it("render the information of all columns", () => {
    const wrapper = shallow(<TableRow columns={columns} row={row} index={0} />);

    columns.forEach((c, i) => {
      expect(wrapper.find("tr").childAt(i)).toHaveText(row[c.key]);
    });
  });
});
