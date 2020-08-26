import * as React from "react";
import AceEditor from "react-ace";

import "./AppValues.css";

interface IAppValuesProps {
  values: string;
}

function AppValues(props: IAppValuesProps) {
  let values = <p>The current application was installed without specifying any values</p>;
  if (props.values !== "") {
    values = (
      <AceEditor
        mode="yaml"
        theme="xcode"
        name="values"
        className="installation-values"
        width="100%"
        maxLines={40}
        setOptions={{ showPrintMargin: false }}
        editorProps={{ $blockScrolling: Infinity }}
        value={props.values}
        readOnly={true}
      />
    );
  }
  return (
    <section aria-labelledby="installation-values">
      <h5 className="section-title" id="installation-values">
        Installation Values
      </h5>
      {values}
    </section>
  );
}

export default AppValues;
