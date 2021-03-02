import React from "react";
import { DataGrid } from "@material-ui/data-grid";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import SearchIcon from "@material-ui/icons/Search";
import axios from "axios";
import UserAdder from "./UserAdder";
import UserDeleter from "./UserDeleter";

const columns = [
  { field: "id", headerName: "번호", width: 90 },
  { field: "email", headerName: "이메일", width: 200 },
  { field: "is_super", headerName: "관리자", width: 100 },
];

const UserlistTemplate = (props) => {
  const { classes, email, setUsers } = props;
  const [rows, setRows] = React.useState([]);
  const [view, setView] = React.useState([]);
  const [selection, setSelection] = React.useState([]);
  React.useEffect(() => {
    const request = {
      method: "get",
      url: "http://13.124.180.188:8000/api/users",
      headers: {
        "X-User-Email": email,
      },
    };
    console.log("[유저목록조회]");
    axios(request)
      .then((response) => {
        var _rows = [];
        const users_count = response.data.data.users_count;
        for (var i = 0; i < users_count; i++) {
          _rows.push({
            id: response.data.data.users[i].user_id,
            email: response.data.data.users[i].user_email,
            is_super: response.data.data.users[i].is_super ? "Yes" : "No",
          });
        }
        console.log(_rows);
        setRows(_rows);
        setView(_rows);
        setUsers(_rows.map((user) => user.email));
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const addUser = (user) => {
    setRows(rows.concat(user), setUsers(rows.map((user) => user.email)));
  };

  const handleChangeSelection = (select) => {
    setSelection(
      rows.filter((row) => select.rowIds.includes(row["id"].toString()))
    );
  };

  const handleChangeSearch = (e) => {
    setView(
      rows.filter((row) => row.email.split("@")[0].includes(e.target.value))
    );
  };
  return (
    <>
      <AppBar
        className={classes.searchBar}
        position="static"
        color="default"
        elevation={0}
      >
        <Toolbar>
          <Grid container spacing={2} alignItems="center">
            <Grid item>
              <SearchIcon className={classes.block} color="inherit" />
            </Grid>
            <Grid item xs>
              <TextField
                fullWidth
                placeholder="이메일 주소로 검색"
                InputProps={{
                  disableUnderline: true,
                  className: classes.searchInput,
                }}
                onChange={handleChangeSearch}
              />
            </Grid>
            <Grid item>
              <UserAdder classes={classes} addUser={addUser} email={email} />
              <UserDeleter
                classes={classes}
                email={email}
                selection={selection}
              />
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      <div style={{ height: 400, width: "100%" }}>
        <DataGrid
          rows={view}
          columns={columns}
          pageSize={5}
          checkboxSelection
          onSelectionChange={handleChangeSelection}
        />
      </div>
    </>
  );
};

export default UserlistTemplate;
