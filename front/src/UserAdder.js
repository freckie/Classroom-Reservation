import React from "react";
import Button from "@material-ui/core/Button";
import AddIcon from "@material-ui/icons/Add";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import Checkbox from "@material-ui/core/Checkbox";
import TextField from "@material-ui/core/TextField";
import Slide from "@material-ui/core/Slide";
import axios from "axios";

const UserAdder = (props) => {
  const { classes, email, addUser } = props;
  const statusPool = ["200", "400", "401", "403", "500", "F01", "F02", "F03"];
  const [status, setStatus] = React.useState();
  const [formOpen, setFormOpen] = React.useState(false);
  const [resultOpen, setResultOpen] = React.useState(false);
  const [inputEmail, setInputEmail] = React.useState("");
  const [inputSuper, setInputSuper] = React.useState(false);
  const handleClickAddBtn = () => {
    setFormOpen(true);
  };

  const handleClickSubmit = () => {
    if (!inputEmail.includes("@") || inputEmail.split("@")[0] === "") {
      // 이메일 형식이 올바르지않음
      setStatus("F01");
    } else if (!inputEmail.endsWith("@khu.ac.kr")) {
      // 경희대학교 이메일 도메인에 속하지 않음
      setStatus("F02");
    } else {
      const request = {
        method: "post",
        url: "http://13.124.180.188:8000/api/users",
        headers: {
          "X-User-Email": email,
        },
        data: {
          email: inputEmail,
          is_super: inputSuper,
        },
      };
      console.log("[유저 등록]");
      axios(request)
        .then((response) => {
          setStatus(response.status.toString());
          if (response.status === 200) {
            addUser({
              id: response.data.data["user_id"],
              email: inputEmail,
              is_super: inputSuper ? "Yes" : "No",
            });
          }
          setResultOpen(true);
        })
        .catch((error) => {
          setStatus("F03");
          setResultOpen(true);
        });
    }
  };

  const handleCloseForm = () => {
    setInputEmail("");
    setInputSuper(false);
    setFormOpen(false);
  };

  const handleCloseResult = () => {
    setResultOpen(false);
    if (status === "200") {
      handleCloseForm();
    }
  };

  return (
    <>
      <Button
        variant="contained"
        color="primary"
        className={classes.addUser}
        onClick={handleClickAddBtn}
      >
        <AddIcon />
      </Button>
      <Dialog open={formOpen} onClose={handleCloseForm} aria-labelledby="form">
        <DialogTitle id="form">유저 등록</DialogTitle>
        <DialogContent>
          <DialogContentText>
            해당 서비스에 접근가능한 유저를 추가합니다.
          </DialogContentText>
          <TextField
            value={inputEmail}
            autoFocus
            margin="dense"
            id="email"
            label="이메일 주소"
            type="email"
            fullWidth="true"
            onInput={(e) => setInputEmail(e.target.value)}
          />
          <p>
            <Checkbox
              checked={inputSuper}
              onChange={(e) => setInputSuper(e.target.checked)}
              id="is_super"
              color="primary"
            />
            관리자 권한 부여
          </p>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseForm}>취소</Button>
          <Button onClick={handleClickSubmit} color="primary">
            등록
          </Button>
        </DialogActions>
        <Dialog
          open={resultOpen}
          onClose={handleCloseResult}
          aria-labelledby="confirm"
        >
          <DialogTitle id="confirm">
            {status === "200" ? "성공" : "에러"}
          </DialogTitle>
          <DialogContent>
            {statusPool.includes(status) ? (
              {
                200: <DialogContent>성공적으로 등록되었습니다!</DialogContent>,
                400: <DialogContent>요청 정보가 부족합니다.</DialogContent>,
                401: (
                  <DialogContent>로그인 정보를 알 수 없습니다.</DialogContent>
                ),
                403: <DialogContent>등록 권한이 없습니다.</DialogContent>,
                500: <DialogContent>에러가 발생하였습니다.</DialogContent>,
                F01: (
                  <DialogContent>
                    이메일 형식이 올바르지 않습니다.
                  </DialogContent>
                ),
                F02: (
                  <DialogContent>
                    경희대학교 계정 이메일만 등록할 수 있습니다.
                  </DialogContent>
                ),
                F03: (
                  <DialogContent>요청 중 에러가 발생하였습니다.</DialogContent>
                ),
              }[status]
            ) : (
              <DialogContent>
                에러가 발생하였습니다. Status[{status}]
              </DialogContent>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseResult}>확인</Button>
          </DialogActions>
        </Dialog>
      </Dialog>
    </>
  );
};

export default UserAdder;
