import React from "react";
import Button from "@material-ui/core/Button";
import AddIcon from "@material-ui/icons/Add";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import axios from "axios";

const FileAdder = (props) => {
  const { classes, email, updateRows, users } = props;
  const [file, setFile] = React.useState("");
  const statusPool = ["200", "400", "401", "403", "404", "500", "F03"];
  const [addOpen, setAddOpen] = React.useState(false);
  const [alertOpen, setAlertOpen] = React.useState(false);
  const [status, setStatus] = React.useState();
  const handleClickAdd = () => {
    setAddOpen(true);
  };

  const handleCloseAdd = () => {
    setFile("");
    setAddOpen(false);
  };

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleClickSubmit = () => {
    console.log("진입1");
    if (file === "") {
      setStatus("400");
      setAlertOpen(true);
    } else {
      console.log("진입2");
      var frm = new FormData();
      frm.append("file", file);
      const request = {
        method: "post",
        url: "http://13.124.180.188:8000/api/files",
        headers: {
          "X-User-Email": email,
          "Content-Type": "multipart/form-data",
        },
        data: frm,
      };
      console.log("[파일 추가]");
      axios(request)
        .then((response) => {
          if (response.status === 200) {
            const request = {
              method: "post",
              url:
                "http://13.124.180.188:8000/api/files/" +
                response.data.file_id +
                "/share",
              headers: {
                "X-User-Email": email,
              },
              data: {
                user_emails: users,
              },
            };
            console.log("[파일 공유]");
            axios(request)
              .then((response) => {
                if (response.status === 200) {
                  const request = {
                    method: "post",
                    url: response.config.url.split("share")[0] + "/protect",
                    headers: {
                      "X-User-Email": email,
                    },
                  };
                  console.log("[파일 보호]");
                  axios(request)
                    .then((response) => {
                      if (response.status === 200) {
                        updateRows();
                      }
                      setStatus(response.status.toString());
                      setAlertOpen(true);
                    })
                    .catch((error) => {
                      setStatus("F03");
                      setAlertOpen(true);
                    });
                } else {
                  setStatus(response.status.toString());
                  setAlertOpen(true);
                }
              })
              .catch((error) => {
                setStatus("F03");
                setAlertOpen(true);
              });
          } else {
            setStatus(response.status.toString());
            setAlertOpen(true);
          }
        })
        .catch((error) => {
          setStatus("F03");
          setAlertOpen(true);
        });
    }
  };

  const handleCloseAlert = () => {
    setAlertOpen(false);
    if (status === "200") {
      handleCloseAdd();
    }
  };
  return (
    <>
      <Button
        variant="contained"
        color="primary"
        className={classes.addUser}
        onClick={handleClickAdd}
      >
        <AddIcon />
      </Button>
      <Dialog
        open={addOpen}
        onClose={handleCloseAdd}
        aria-labelledby="form-dialog-title"
      >
        <DialogTitle id="form-dialog-title">파일 추가</DialogTitle>
        <DialogContent>
          <DialogContentText>
            시간표 파일을 추가하고 유저들에게 공유합니다.
          </DialogContentText>
          <div id="upload-box">
            <label htmlFor="raised-button-file">
              <input
                accept="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/vnd.ms-excel"
                className={classes.input}
                style={{ display: "none" }}
                id="raised-button-file"
                type="file"
                onChange={handleFileChange}
              />
              <Button color="primary" variant="contained" component="span">
                Upload
              </Button>
              <p>{file.name}</p>
            </label>
          </div>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseAdd}>취소</Button>
          <Button onClick={handleClickSubmit} color="primary">
            추가
          </Button>
          <Dialog
            open={alertOpen}
            onClose={handleCloseAlert}
            aria-labelledby="alert"
          >
            <DialogTitle id="alert">
              {status === "200" ? "성공" : "에러"}
            </DialogTitle>
            <DialogContent>
              {statusPool.includes(status) ? (
                {
                  200: (
                    <DialogContent>성공적으로 등록되었습니다!</DialogContent>
                  ),
                  400: <DialogContent>파일을 업로드 해주세요.</DialogContent>,
                  401: (
                    <DialogContent>로그인 정보를 알 수 없습니다.</DialogContent>
                  ),
                  403: <DialogContent>권한이 없습니다.</DialogContent>,
                  404: <DialogContent>해당 파일이 없습니다.</DialogContent>,
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
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseAlert}>확인</Button>
            </DialogActions>
          </Dialog>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default FileAdder;
