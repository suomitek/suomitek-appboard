import React, { useState } from "react";

import actions from "actions";
import { CdsButton, CdsIcon } from "components/Clarity/clarity";
import { useDispatch, useSelector } from "react-redux";
import { IStoreState } from "shared/types";

export function AppRepoRefreshAllButton() {
  const [refreshing, setRefreshing] = useState(false);
  const { repos } = useSelector((state: IStoreState) => state.repos);
  const dispatch = useDispatch();

  const handleResyncAllClick = async () => {
    // Fake timeout to show progress
    // TODO(andresmgot): Ideally, we should show the progress of the sync but we don't
    // have that info yet: https://github.com/suomitek/suomitek-appboard/issues/153
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 500);
    if (repos) {
      const repoObjects = repos.map(repo => {
        return {
          name: repo.metadata.name,
          namespace: repo.metadata.namespace,
        };
      });
      dispatch(actions.repos.resyncAllRepos(repoObjects));
    }
  };
  return (
    <div className="refresh-all-button">
      <CdsButton action="outline" onClick={handleResyncAllClick} disabled={refreshing}>
        <CdsIcon shape="refresh" inverse={true} /> {refreshing ? "Refreshing" : "Refresh All"}
      </CdsButton>
    </div>
  );
}
