import actions from "actions";
import { CdsButton } from "components/Clarity/clarity";
import { shallow } from "enzyme";
import * as React from "react";
import { act } from "react-dom/test-utils";
import * as ReactRedux from "react-redux";
import { ISecret } from "shared/types";
import AppRepoAddDockerCreds from "./AppRepoAddDockerCreds.v2";

const secret1 = {
  metadata: {
    name: "foo",
  },
} as ISecret;
const secret2 = {
  metadata: {
    name: "bar",
  },
} as ISecret;
const defaultProps = {
  imagePullSecrets: [],
  togglePullSecret: jest.fn(),
  selectedImagePullSecrets: {},
  namespace: "default",
};

let spyOnUseDispatch: jest.SpyInstance;
const kubeaActions = { ...actions.kube };
beforeEach(() => {
  actions.repos = {
    ...actions.repos,
    createDockerRegistrySecret: jest.fn(),
  };
  const mockDispatch = jest.fn(r => r);
  spyOnUseDispatch = jest.spyOn(ReactRedux, "useDispatch").mockReturnValue(mockDispatch);
});

afterEach(() => {
  actions.kube = { ...kubeaActions };
  spyOnUseDispatch.mockRestore();
});

it("shows an info message if there are no secrets", () => {
  const wrapper = shallow(<AppRepoAddDockerCreds {...defaultProps} />);
  expect(wrapper.text()).toContain("No existing credentials found");
});

it("shows the list of available pull secrets", () => {
  const wrapper = shallow(
    <AppRepoAddDockerCreds {...defaultProps} imagePullSecrets={[secret1, secret2]} />,
  );
  expect(wrapper.text()).toContain(secret1.metadata.name);
  expect(wrapper.text()).toContain(secret2.metadata.name);
});

it("select secrets", () => {
  const wrapper = shallow(
    <AppRepoAddDockerCreds
      {...defaultProps}
      imagePullSecrets={[secret1, secret2]}
      selectedImagePullSecrets={{ [secret1.metadata.name]: true }}
    />,
  );
  const totalCheckbox = wrapper.find("input").filterWhere(i => i.prop("type") === "checkbox");
  expect(totalCheckbox.length).toBe(2);

  const selectedCheckbox = totalCheckbox.filterWhere(i => i.prop("checked") === true);
  expect(selectedCheckbox.length).toBe(1);
});

it("renders the form to create a registry secret", () => {
  const wrapper = shallow(<AppRepoAddDockerCreds {...defaultProps} />);

  expect(wrapper.text()).not.toContain("Secret Name");

  const button = wrapper.find(CdsButton).filterWhere(b => b.html().includes("Add new"));
  act(() => {
    (button.prop("onClick") as any)();
  });
  wrapper.update();

  expect(wrapper.text()).toContain("Secret Name");
});

it("submits the new secret", async () => {
  const createDockerRegistrySecret = jest.fn().mockReturnValue(true);
  actions.repos = {
    ...actions.repos,
    createDockerRegistrySecret,
  };
  const wrapper = shallow(<AppRepoAddDockerCreds {...defaultProps} />);
  // Open form
  const button = wrapper.find(CdsButton).filterWhere(b => b.html().includes("Add new"));
  act(() => {
    (button.prop("onClick") as any)();
  });
  wrapper.update();

  const secretName = "repo-1";
  const user = "foo";
  const password = "pass";
  const email = "foo@bar.com";
  const server = "docker.io";

  wrapper
    .find("#kubeapps-docker-cred-secret-name")
    .simulate("change", { target: { value: secretName } });
  wrapper.find("#kubeapps-docker-cred-server").simulate("change", { target: { value: server } });
  wrapper.find("#kubeapps-docker-cred-username").simulate("change", { target: { value: user } });
  wrapper
    .find("#kubeapps-docker-cred-password")
    .simulate("change", { target: { value: password } });
  wrapper.find("#kubeapps-docker-cred-email").simulate("change", { target: { value: email } });
  wrapper.update();

  const submit = wrapper.find(CdsButton).filterWhere(b => b.html().includes("Submit"));
  await act(async () => {
    await (submit.prop("onClick") as () => Promise<any>)();
  });
  wrapper.update();

  expect(createDockerRegistrySecret).toHaveBeenCalledWith(
    secretName,
    user,
    password,
    email,
    server,
    defaultProps.namespace,
  );
  // There should be a new item with the secret
  expect(wrapper.find("#app-repo-secret-repo-1")).toExist();
});
