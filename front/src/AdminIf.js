import React from "react";
import Paperbase from "./Paperbase";
import Forbidden from "./Forbidden";

function AdminIf(props) {
  const { email, ...rest } = props;
  return email !== "" ? <Paperbase email={email} {...rest} /> : <Forbidden />;
}

export default AdminIf;
