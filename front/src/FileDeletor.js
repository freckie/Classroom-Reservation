import React from "react";
import Button from "@material-ui/core/Button";
import RemoveIcon from "@material-ui/icons/Remove";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import axios from "axios";

const FileDeletor = (props) => {
  const { classes, email, updateRows, selection } = props;
  const [deleteOpen, setDeleteOpen] = React.useState(false);
  const [alertOpen, setAlertOpen] = React.useState(false);
  const [status, setStatus] = React.useState();
  const statusPool = ["200", "401", "403", "500", "F03"];
  const handleClickDelete = () => {
    setDeleteOpen(true);
  };
  const handleCloseDelete = () => {
    setDeleteOpen(false);
  };

  const handleClickDeletion = () => {
    setStatus("200");
    console.log("[파일삭제]");
    console.log(selection.length);
    var statuses = Array.from({ length: selection.length }, () => "");
    function checkStatuses() {
      console.log(statuses);
      for (var j = 0; j < statuses.length; ++j) {
        if (statuses[j] === "") {
          return;
        } else if (statuses[j] !== "200") {
          setStatus(statuses[j]);
        }
      }
      updateRows();
      setAlertOpen(true);
    }
    for (var i = 0; i < selection.length; ++i) {
      console.log(selection[i].file_id);
      const request = {
        method: "delete",
        url: "http://13.124.180.188:8000/api/files/" + selection[i].file_id,
        headers: {
          "X-User-Email": email,
        },
      };

      axios(request)
        .then((response) => {
          console.log(response);
          statuses[i] = response.status.toString();
          checkStatuses();
        })
        .catch((error) => {
          console.log(error);
          statuses[i] = "F03";
          checkStatuses();
        });
    }
  };

  const handleCloseAlert = () => {
    setAlertOpen(false);
    handleCloseDelete();
  };

  const handleClickAlert = () => {
    handleCloseAlert();
  };

  return (
    <>
      <Button
        variant="contained"
        color="secondary"
        className={classes.removefile}
        onClick={handleClickDelete}
      >
        <RemoveIcon />
      </Button>
      <Dialog
        open={deleteOpen}
        onClose={handleCloseDelete}
        scroll="paper"
        aria-labelledby="scroll-dialog-title"
        aria-describedby="scroll-dialog-description"
      >
        <DialogTitle id="scroll-dialog-title">파일 삭제</DialogTitle>
        <DialogContent dividers="true">
          <DialogContentText id="scroll-dialog-description">
            {selection.length === 0 ? (
              "선택된 파일이 없습니다."
            ) : (
              <>
                <p>다음 파일이 완전히 삭제됩니다.</p>
                <List dense="true">
                  {selection.map((x) => (
                    <ListItem>
                      <ListItemText
                        primary={x.file_name}
                        secondary={x.created_at + " 생성"}
                      ></ListItemText>
                    </ListItem>
                  ))}
                </List>
              </>
            )}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          {selection.length === 0 ? (
            <Button onClick={handleCloseDelete}>확인</Button>
          ) : (
            <>
              <Button onClick={handleCloseDelete}>취소</Button>
              <Button onClick={handleClickDeletion} color="secondary">
                삭제
              </Button>
              <Dialog open={alertOpen} onClose={handleCloseAlert}>
                {statusPool.includes(status) ? (
                  {
                    200: (
                      <DialogContent>성공적으로 삭제되었습니다.</DialogContent>
                    ),
                    401: (
                      <DialogContent>
                        로그인 정보를 알 수 없습니다.
                      </DialogContent>
                    ),
                    403: <DialogContent>권한이 없습니다.</DialogContent>,
                    500: <DialogContent>에러가 발생하였습니다.</DialogContent>,
                    F03: (
                      <DialogContent>
                        요청 중 에러가 발생하였습니다.
                      </DialogContent>
                    ),
                  }[status]
                ) : (
                  <DialogContent>
                    에러가 발생하였습니다. Status[{status}]
                  </DialogContent>
                )}
                <DialogActions>
                  <Button onClick={handleClickAlert}>확인</Button>
                </DialogActions>
              </Dialog>
            </>
          )}
        </DialogActions>
      </Dialog>
    </>
  );
};

export default FileDeletor;
