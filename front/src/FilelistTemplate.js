import React from "react";
import { DataGrid } from "@material-ui/data-grid";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Grid from "@material-ui/core/Grid";
import Button from "@material-ui/core/Button";
import TextField from "@material-ui/core/TextField";
import SearchIcon from "@material-ui/icons/Search";
import RemoveIcon from "@material-ui/icons/Remove";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import axios from "axios";
import FileAdder from "./FileAdder";
import FileDeletor from "./FileDeletor";
const columns = [
  { field: "id", headerName: "번호", width: 90 },
  { field: "file_name", headerName: "파일 이름", width: 300 },
  { field: "created_at", headerName: "생성된 시간", width: 200 },
];

const FilelistTemplate = (props) => {
  const { classes, email, users } = props;
  const [rows, setRows] = React.useState([]);
  const [view, setView] = React.useState([]);
  const [selection, setSelection] = React.useState([]);
  const [ai, setAi] = React.useState(0);

  const updateRows = () => {
    const request = {
      method: "get",
      url: "http://13.124.180.188:8000/api/files",
      headers: {
        "X-User-Email": email,
      },
    };
    console.log("[파일목록조회]");
    axios(request)
      .then((response) => {
        var _rows = [];
        const files_count = response.data.data.files_count;
        for (var i = 0; i < files_count; i++) {
          _rows.push({
            id: i + 1,
            file_id: response.data.data.files[i].file_id,
            file_name: response.data.data.files[i].file_name,
            created_at: response.data.data.files[i].created_at,
          });
        }
        setAi(files_count + 1);
        setRows(_rows);
        setView(_rows);
      })
      .catch((error) => {
        console.log(error);
      });
  };
  React.useEffect(() => {
    updateRows();
  }, []);

  const handleChangeSelection = (select) => {
    console.log(select);
    setSelection(
      rows.filter((row) => select.rowIds.includes(row["id"].toString()))
    );
  };

  const handleChangeSearch = (e) => {
    setView(rows.filter((row) => row.file_name.includes(e.target.value)));
  };

  const [deleteOpen, setDeleteOpen] = React.useState(false);

  const handleClickDelete = () => {
    setDeleteOpen(true);
  };

  const handleCloseDelete = () => {
    setDeleteOpen(false);
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
                placeholder="파일 이름으로 검색"
                InputProps={{
                  disableUnderline: true,
                  className: classes.searchInput,
                }}
                onChange={handleChangeSearch}
              />
            </Grid>
            <Grid item>
              <FileAdder
                classes={classes}
                updateRows={updateRows}
                email={email}
                users={users}
              />
              <FileDeletor
                classes={classes}
                email={email}
                updateRows={updateRows}
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

export default FilelistTemplate;
