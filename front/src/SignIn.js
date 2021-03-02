import React from "react";
import Avatar from "@material-ui/core/Avatar";
import Button from "@material-ui/core/Button";
import CssBaseline from "@material-ui/core/CssBaseline";
import TextField from "@material-ui/core/TextField";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import Checkbox from "@material-ui/core/Checkbox";
import Link from "@material-ui/core/Link";
import Grid from "@material-ui/core/Grid";
import Box from "@material-ui/core/Box";
import LockOutlinedIcon from "@material-ui/icons/LockOutlined";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import Container from "@material-ui/core/Container";
import { Redirect } from "react-router";
import axios from "axios";
import FileSelectDialog from "./FileSelectDialog";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {"Copyright © "}
      <Link color="inherit" href="http://swedu.khu.ac.kr/">
        경희대학교 SW중심대학사업단
      </Link>{" "}
      {new Date().getFullYear()}
      {"."}
    </Typography>
  );
}

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main,
  },
  form: {
    width: "100%", // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));

export default function SignIn({ setEmail, history, location }) {
  const classes = useStyles();
  const [users, setUsers] = React.useState([]);
  const [input, setInput] = React.useState("");
  const [open, setOpen] = React.useState(false);
  const [alert, setAlert] = React.useState(false);
  React.useEffect(() => {
    // 차후, 인증구조 반드시 대체 필요!
    setEmail("");
    const request = {
      method: "get",
      url: "http://13.124.180.188:8000/api/users",
      headers: {
        "X-User-Email": "super@khu.ac.kr",
      },
    };
    console.log("[유저목록조회]");
    axios(request)
      .then((response) => {
        setUsers(
          response.data.data.users.map((x) => {
            const user = {
              id: x.user_id,
              email: x.user_email,
              is_super: x.is_super,
            };

            return user;
          })
        );
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);
  const handleChangeEmail = (e) => {
    setInput(e.target.value);
  };
  const handleClickLogin = () => {
    const user = users.find((user) => user.email === input);
    if (!user) {
      setAlert(true);
    } else {
      console.log(user);
      setEmail(input);
      if (user.is_super) {
        history.push("/admin");
      } else {
        setOpen(true);
      }
    }
  };
  const handleCloseAlert = () => {
    setAlert(false);
  };
  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Typography component="h5" variant="h5">
          경희대학교 대면 시험 강의실 예약
        </Typography>

        <TextField
          variant="outlined"
          margin="normal"
          required
          fullWidth
          id="email"
          label="경희대학교 이메일"
          name="email"
          type="email"
          autoFocus
          onChange={handleChangeEmail}
        />
        <Button
          type="submit"
          fullWidth
          variant="contained"
          color="primary"
          className={classes.submit}
          onClick={handleClickLogin}
        >
          로그인
        </Button>
        <Dialog
          open={alert}
          onClose={handleCloseAlert}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">로그인 실패</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              등록된 유저 이메일을 찾을 수 없습니다.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseAlert} color="primary" autoFocus>
              확인
            </Button>
          </DialogActions>
        </Dialog>
        <FileSelectDialog open={open} setOpen={setOpen} />
      </div>
      <Box mt={8}>
        <Copyright />
      </Box>
    </Container>
  );
}
