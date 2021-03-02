import React from "react";
import PropTypes from "prop-types";
import Paper from "@material-ui/core/Paper";
import { withStyles } from "@material-ui/core/styles";
import UserlistTemplate from "./UserlistTemplate";
import FilelistTemplate from "./FilelistTemplate";
const styles = (theme) => ({
  paper: {
    maxWidth: 936,
    margin: "auto",
    overflow: "hidden",
  },
  searchBar: {
    borderBottom: "1px solid rgba(0, 0, 0, 0.12)",
  },
  searchInput: {
    fontSize: theme.typography.fontSize,
  },
  block: {
    display: "block",
  },
  addUser: {
    marginRight: theme.spacing(1),
  },

  contentWrapper: {
    margin: "40px 16px",
  },
});

function Content(props) {
  const { classes, category, email} = props;
  const [users, setUsers] = React.useState([]);
  const switchCategoryTemplate = () => {
    switch (category) {
      case "유저관리":
        return (
          <UserlistTemplate
            classes={classes}
            email={email}
            setUsers={setUsers}
          />
        );
        break;
      case "파일관리":
        return (
          <FilelistTemplate classes={classes} email={email} users={users} />
        );
        break;
    }
  };
  return <Paper className={classes.paper}>{switchCategoryTemplate()}</Paper>;
}

Content.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Content);
