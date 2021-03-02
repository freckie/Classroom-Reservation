import React from "react";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import { makeStyles } from "@material-ui/core/styles";
import ListItemAvatar from "@material-ui/core/ListItemAvatar";
import DescriptionIcon from "@material-ui/icons/Description";
import ListItemText from "@material-ui/core/ListItemText";
import Avatar from "@material-ui/core/Avatar";
import axios from "axios";
import DialogTitle from "@material-ui/core/DialogTitle";
import Dialog from "@material-ui/core/Dialog";
import { blue } from "@material-ui/core/colors";

const useStyles = makeStyles({
  avatar: {
    backgroundColor: blue[100],
    color: blue[600],
  },
});

function FileSelectDialog(props) {
  const classes = useStyles();
  const { open, setOpen } = props;
  const [files, setFiles] = React.useState([]);
  React.useEffect(() => {
    const request = {
      method: "get",
      url: "http://13.124.180.188:8000/api/files",
      headers: {
        "X-User-Email": "super@khu.ac.kr",
      },
    };
    console.log("[파일목록조회]");
    axios(request)
      .then((response) => {
        if (response.data.data.files_count !== 0) {
          setFiles(response.data.data.files);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);
  const handleClose = () => {
    setOpen(false);
  };
  const handleListItemClick = (file) => {
    window.location.href =
      "https://docs.google.com/spreadsheets/d/" + file.file_id + "/edit";
  };
  return (
    <Dialog
      onClose={handleClose}
      aria-labelledby="simple-dialog-title"
      open={open}
    >
      <DialogTitle id="simple-dialog-title">파일선택</DialogTitle>
      <List>
        {files.map((file) => (
          <ListItem
            button
            onClick={() => handleListItemClick(file)}
            key={file.file_id}
          >
            <ListItemAvatar>
              <Avatar className={classes.avatar}>
                <DescriptionIcon />
              </Avatar>
            </ListItemAvatar>
            <ListItemText primary={file.file_name} />
          </ListItem>
        ))}
      </List>
    </Dialog>
  );
}

export default FileSelectDialog;
