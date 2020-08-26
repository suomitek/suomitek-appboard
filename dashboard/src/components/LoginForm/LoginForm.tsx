import { Location } from "history";
import * as React from "react";
import { Lock } from "react-feather";
import { Redirect } from "react-router";

import LoadingWrapper from "../../components/LoadingWrapper";
import "./LoginForm.css";

export interface ILoginFormProps {
  cluster: string;
  authenticated: boolean;
  authenticating: boolean;
  authenticationError: string | undefined;
  oauthLoginURI: string;
  authenticate: (cluster: string, token: string) => any;
  checkCookieAuthentication: (cluster: string) => void;
  location: Location;
}

interface ILoginFormState {
  token: string;
}

class LoginForm extends React.Component<ILoginFormProps, ILoginFormState> {
  public state: ILoginFormState = { token: "" };

  public componentDidMount() {
    if (this.props.oauthLoginURI) {
      this.props.checkCookieAuthentication(this.props.cluster);
    }
  }

  public render() {
    if (this.props.authenticating) {
      return <LoadingWrapper />;
    }
    if (this.props.authenticated) {
      const { from } = (this.props.location.state as any) || { from: { pathname: "/" } };
      return <Redirect to={from} />;
    }

    return (
      <section className="LoginForm">
        <div className="LoginForm__container padding-v-bigger bg-skew">
          <div className="container container-tiny">
            {this.props.authenticationError && (
              <div className="alert alert-error margin-c" role="alert">
                There was an error connecting to the Kubernetes API. Please check that your token is
                valid.
              </div>
            )}
          </div>
          <div className="bg-skew__pattern bg-skew__pattern-dark type-color-reverse">
            <div className="container">
              <h2>
                <Lock /> Login
              </h2>
              {this.props.oauthLoginURI ? this.oauthLogin() : this.tokenLogin()}
            </div>
          </div>
        </div>
      </section>
    );
  }

  private handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const { token } = this.state;
    return token && (await this.props.authenticate(this.props.cluster, token));
  };

  private handleTokenChange = (e: React.FormEvent<HTMLInputElement>) => {
    this.setState({ token: e.currentTarget.value });
  };

  private oauthLogin = () => {
    return (
      <>
        <p>
          Your cluster operator has enabled login via an authentication provider.{" "}
          <a
            href="https://github.com/kubeapps/kubeapps/blob/master/docs/user/using-an-OIDC-provider.md"
            target="_blank"
            rel="noopener noreferrer"
          >
            Click here
          </a>{" "}
          for more info about using authentication providers with Kubeapps.
        </p>
        <a href={this.props.oauthLoginURI} className="button button-accent">
          Login
        </a>
      </>
    );
  };

  private tokenLogin = () => {
    return (
      <>
        <p>
          Your cluster operator should provide you with a Kubernetes API token.{" "}
          <a
            href="https://github.com/kubeapps/kubeapps/blob/master/docs/user/access-control.md"
            target="_blank"
            rel="noopener noreferrer"
          >
            Click here
          </a>{" "}
          for more info on how to create and use Bearer Tokens.
        </p>
        <div className="bg-skew__content">
          <form onSubmit={this.handleSubmit}>
            <div>
              <label htmlFor="token">Kubernetes API Token</label>
              <input
                name="token"
                id="token"
                type="password"
                placeholder="Token"
                required={true}
                onChange={this.handleTokenChange}
                value={this.state.token}
              />
            </div>
            <p>
              <button type="submit" className="button button-accent">
                Login
              </button>
            </p>
          </form>
        </div>
      </>
    );
  };
}

export default LoginForm;
