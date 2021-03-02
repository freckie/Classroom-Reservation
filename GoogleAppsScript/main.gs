/**
 * Const Values
 */

var host = ;
var port = ;
var ErrorCodeMessages = {
  undefined: "등록되지 않은 오류가 발생했습니다. 다시 시도해주세요.",
  "error-00": "셀을 다시 선택해주세요.",
  "error-01": "한번에 한 강의실만 선택해주세요.",
  "error-02": "상단 혹은 좌측의 Header는 선택에서 제외해주세요.",
  "error-03": "해당 강의실의 수용인원이 초과되었습니다.",
  "error-04": "선택한 시간 범위가 잘못되었습니다.",
  "error-05": "이미 예약된 시간입니다.",
};

/**
 * Script Functions
 */

/**
 * Create custom menu and show sidebar when the sheet is opened.
 *
 * @param {Object} e the event parameter for a simple onOpen trigger.
 */
function onOpen(e) {
  var ui = SpreadsheetApp.getUi();
  SpreadsheetApp.getUi()
    .createAddonMenu()
    .addItem("스크립트 재구동", "showSidebar")
    .addToUi();
}

/**
 * Runs when the add-on is installed; calls onOpen() to ensure menu creation and
 * any other initializion work is done immediately.
 *
 * @param {Object} e The event parameter for a simple onInstall trigger.
 */
function onInstall(e) {
  onOpen(e);
}

/**
 * Opens a sidebar. The sidebar structure is described in the Sidebar.html
 * project file.
 */
function showSidebar() {
  var ui = HtmlService.createTemplateFromFile("Sidebar")
    .evaluate()
    .setSandboxMode(HtmlService.SandboxMode.IFRAME)
    .setTitle("KHU Classroom Lender");
  SpreadsheetApp.getUi().showSidebar(ui);
}

/**
 * Get range of selected cells.
 *
 * @return {values, row, col, width, height}: values and xywh coordinates
 * @return topHeader: header value
 * @return date: date value
 */

function requestReservationInfo(range, topHeader, date, column, start, end) {
  var ui = SpreadsheetApp.getUi();

  var fileId = SpreadsheetApp.getActiveSpreadsheet().getId();
  var sheetId = SpreadsheetApp.getActiveSheet().getSheetId();

  var uri = "/api/files/" + fileId + "/" + sheetId + "/cell";
  var queryParam = "?" + "column=" + column + "&start=" + start + "&end=" + end;
  var url = "http://" + host + ":" + port + uri + queryParam;
  var email = Session.getActiveUser().getEmail();
  var options = {
    method: "get",
    headers: {
      "X-User-Email": email,
    },
    muteHttpExceptions: true,
  };
  var response = UrlFetchApp.fetch(url, options);
  var status = response.getResponseCode();
  if (status === 200) {
    var data = JSON.parse(response.getContentText()).data;
    return [range, topHeader, date, data, email]; // Date type 때문에 에러뜰수도??
  } else if (status === 403) {
    ui.alert("[error 예약정보 받기]:403");
  } else if (status === 404) {
    ui.alert("[error 예약정보 받기]:404");
  }
}

function getRange() {
  var leftHeader = 2;
  var rightHeader = 35;

  var sheet = SpreadsheetApp.getActiveSheet();
  var range = sheet.getActiveRange();
  var values = range.getValues();
  var row = range.getRow();
  var col = range.getColumn();
  var width = range.getWidth();
  var height = range.getHeight();
  for (var i = 0; i < height; ++i) {
    for (var j = 0; j < width; ++j) {
      if (values[i][j] instanceof Date) {
        values[i][j] = values[i][j].toTimeString();
      }
    }
  }
  var topHeader = "s";
  if (col > leftHeader && col < rightHeader)
    topHeader = sheet.getRange(1, col).getValue();
  var date = sheet.getRange(1, 1).getValue();
  return [{ values, row, col, width, height }, topHeader, date];
}

/**
 * Runs when the submit button in Sidebar is clicked.
 * Merges selected cells and put value(classname, prof. name)
 * If error occurs, opens alert dialog.
 */
function submit(inputValues, errorCode) {
  var ui = SpreadsheetApp.getUi();

  // 에러가 발생했는지 확인
  if (errorCode !== "") {
    // alert 출력
    var msg =
      errorCode in ErrorCodeMessages
        ? ErrorCodeMessages[errorCode]
        : ErrorCodeMessages[undefined];
    ui.alert(msg);
    return;
  }

  // 변수들이 모두 입력되었는지 확인
  if (
    inputValues.professor === "" ||
    inputValues.lecture === "" ||
    inputValues.capacity === ""
  ) {
    ui.alert("누락된 정보가 있습니다. 모두 입력해주세요.");
    return;
  }

  // 최종 확인
  var msg =
    "예약 대상이 [" +
    inputValues.classRoom +
    " / " +
    inputValues.time +
    "] 이 맞습니까?";
  var response = ui.alert("예약 확인", msg, ui.ButtonSet.YES_NO);
  if (response != ui.Button.YES) {
    return;
  }

  var fileId = SpreadsheetApp.getActiveSpreadsheet().getId();
  var sheetId = SpreadsheetApp.getActiveSheet().getSheetId();

  var uri = "/api/files/" + fileId + "/" + sheetId + "/reservation";
  var url = "http://" + host + ":" + port + uri;
  var email = Session.getActiveUser().getEmail();
  var data = {
    column: inputValues.column,
    start: parseInt(inputValues.start),
    end: parseInt(inputValues.end),
    lecture: inputValues.lecture,
    professor: inputValues.professor,
    capacity: parseInt(inputValues.capacity),
  };
  var options = {
    method: "post",
    contentType: "application/json",
    headers: {
      "X-User-Email": email,
    },
    muteHttpExceptions: true,
    payload: JSON.stringify(data),
  };
  var response = UrlFetchApp.fetch(url, options);
  var status = response.getResponseCode();
  Logger.log(response);
  Logger.log(status);

  if (status === 200) {
    ui.alert("예약 생성에 성공했습니다.");
  } else {
    ui.alert("[error 예약하기]:" + status);
  }

  return status;
}

function cancel(id, column, start, end) {
  var ui = SpreadsheetApp.getUi();
  var msg = "해당 예약을 정말 취소하시겠습니까?";
  var response = ui.alert("예약 취소", msg, ui.ButtonSet.YES_NO);
  if (response != ui.Button.YES) {
    return;
  }

  var fileId = SpreadsheetApp.getActiveSpreadsheet().getId();
  var sheetId = SpreadsheetApp.getActiveSheet().getSheetId();

  var uri = "/api/files/" + fileId + "/" + sheetId + "/reservation/" + id;
  var url = "http://" + host + ":" + port + uri;
  var email = Session.getActiveUser().getEmail();
  var options = {
    method: "delete",
    headers: {
      "X-User-Email": email,
    },
    muteHttpExceptions: true,
  };
  var response = UrlFetchApp.fetch(url, options);
  var status = response.getResponseCode();

  // 요청 성공한 경우
  if (status === 200) {
    ui.alert("예약이 취소되었습니다.");
  }
  // 요청 실패한 경우
  else {
    ui.alert("[error 예약취소]:", status);
  }
}
