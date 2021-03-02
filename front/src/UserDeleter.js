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

const UserDeleter = (props) => {
  const [deleteOpen, setDeleteOpen] = React.useState(false);
  const [alertOpen, setAlertOpen] = React.useState(false);

  const handleClickDelete = () => {
    setDeleteOpen(true);
  };
  const handleCloseDelete = () => {
    setDeleteOpen(false);
  };

  const handleClickDeletion = () => {
    // 추후에 유저삭제 요청보내기
    setAlertOpen(true);
  };

  const handleCloseAlert = () => {
    setAlertOpen(false);
    handleCloseDelete();
  };

  const handleClickAlert = () => {
    handleCloseAlert();
  };
  const descriptionElementRef = React.useRef(null);
  React.useEffect(() => {
    if (deleteOpen) {
      const { current: descriptionElement } = descriptionElementRef;
      if (descriptionElement !== null) {
        descriptionElement.focus();
      }
    }
  }, [deleteOpen]);
  const { classes, email, selection } = props;
  return (
    <>
      <Button
        variant="contained"
        color="secondary"
        className={classes.removeUser}
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
        <DialogTitle id="scroll-dialog-title">유저 삭제</DialogTitle>
        <DialogContent dividers="true">
          <DialogContentText
            id="scroll-dialog-description"
            ref={descriptionElementRef}
            tabIndex={-1}
          >
            {selection.length === 0 ? (
              "선택된 유저가 없습니다."
            ) : (
              <>
                <p>다음 유저 데이터가 완전히 삭제됩니다.</p>
                <List dense="true">
                  {selection.map((x) => (
                    <ListItem>
                      <ListItemText
                        primary={x.email}
                        secondary={"관리자 권한 : " + x.is_super}
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
                <DialogContent>추후에 업데이트 예정입니다.</DialogContent>
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

export default UserDeleter;
