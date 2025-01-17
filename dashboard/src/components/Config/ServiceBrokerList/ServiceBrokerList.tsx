import * as React from "react";

import { Link } from "react-router-dom";
import * as url from "shared/url";
import LoadingWrapper from "../../../components/LoadingWrapper";
import PageHeader from "../../../components/PageHeader";
import { IServiceCatalogState } from "../../../reducers/catalog";
import { IServiceBroker } from "../../../shared/ServiceCatalog";
import { IRBACRole } from "../../../shared/types";
import { CardGrid } from "../../Card";
import {
  ErrorSelector,
  MessageAlert,
  ServiceBrokersNotFoundAlert,
  ServiceCatalogNotInstalledAlert,
} from "../../ErrorAlert";
import ServiceBrokerItem from "./ServiceBrokerItem";

interface IServiceBrokerListProps {
  errors: {
    fetch?: Error;
    update?: Error;
  };
  getBrokers: () => Promise<any>;
  brokers: IServiceCatalogState["brokers"];
  sync: (broker: IServiceBroker) => Promise<any>;
  checkCatalogInstalled: () => Promise<any>;
  isInstalled: boolean;
  cluster: string;
}

export const RequiredRBACRoles: { [s: string]: IRBACRole[] } = {
  resync: [
    {
      apiGroup: "servicecatalog.k8s.io",
      clusterWide: true,
      resource: "clusterservicebrokers",
      verbs: ["patch"],
    },
  ],
  view: [
    {
      apiGroup: "servicecatalog.k8s.io",
      clusterWide: true,
      resource: "clusterservicebrokers",
      verbs: ["list"],
    },
  ],
};

class ServiceBrokerList extends React.Component<IServiceBrokerListProps> {
  public componentDidMount() {
    this.props.checkCatalogInstalled();
    this.props.getBrokers();
  }

  public render() {
    const { brokers, errors, sync } = this.props;
    let body = <span />;
    if (errors.fetch) {
      body = (
        <ErrorSelector
          error={errors.fetch}
          resource="Service Brokers"
          action="view"
          defaultRequiredRBACRoles={RequiredRBACRoles}
        />
      );
    } else {
      if (brokers.list.length > 0) {
        if (errors.update) {
          body = (
            <ErrorSelector
              error={errors.update}
              resource="Service Brokers"
              action="resync"
              defaultRequiredRBACRoles={RequiredRBACRoles}
            />
          );
        } else {
          body = (
            <CardGrid className="BrokerList">
              {brokers.list.map(broker => (
                <ServiceBrokerItem key={broker.metadata.uid} broker={broker} sync={sync} />
              ))}
            </CardGrid>
          );
        }
      } else {
        body = <ServiceBrokersNotFoundAlert />;
      }
    }
    return (
      <section className="AppList">
        <PageHeader>
          <h1>Service Brokers</h1>
        </PageHeader>
        <LoadingWrapper loaded={!brokers.isFetching}>{this.renderBody(body)}</LoadingWrapper>
      </section>
    );
  }

  private renderBody(body: React.ReactFragment) {
    const { cluster, isInstalled } = this.props;
    if (cluster !== "default") {
      return (
        <MessageAlert header="Service brokers can be created on the default cluster only">
          <div>
            <p className="margin-v-normal">
              Suomitek-appboard' Service Broker support enables the addition of{" "}
              <Link to={url.app.config.brokers("default")}>
                service brokers on the default cluster only
              </Link>
              .
            </p>
          </div>
        </MessageAlert>
      );
    }

    if (!isInstalled) {
      return <ServiceCatalogNotInstalledAlert />;
    }

    return <main>{body}</main>;
  }
}

export default ServiceBrokerList;
