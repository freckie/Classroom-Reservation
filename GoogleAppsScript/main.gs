/**
 * Const Values
 */

var ErrorCodeMessages = {
  undefined: '등록되지 않은 오류가 발생했습니다. 다시 시도해주세요.',
  'error-00': '셀을 다시 선택해주세요.',
  'error-01': '한번에 한 강의실만 선택해주세요.',
  'error-02': '상단 혹은 좌측의 Header는 선택에서 제외해주세요.',
  'error-03': '해당 강의실의 수용인원이 초과되었습니다.',
  'error-04': '선택한 시간 범위가 잘못되었습니다.',
  'error-05': '이미 예약된 시간입니다.'
};


function onOpen(e) {
  SpreadsheetApp.getUi()
      .createAddonMenu()
      .addItem('View records', 'showSidebar')
      .addToUi();
  showSidebar();
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
  var ui = HtmlService.createTemplateFromFile('Sidebar')
      .evaluate()
      .setSandboxMode(HtmlService.SandboxMode.IFRAME)
      .setTitle('KHU Classroom Lender');
  SpreadsheetApp.getUi().showSidebar(ui);
}

function getRange() {
  var sheet = SpreadsheetApp.getActiveSheet();
  var range = sheet.getActiveRange();
  var values = range.getValues();
  var row = range.getRow();
  var col = range.getColumn();
  var width = range.getWidth();
  var height = range.getHeight();
  var topHeader = sheet.getRange(1, col).getValue();
  var date = sheet.getRange(1, 1).getValue();
  return [{values, row, col, width, height}, topHeader, date];
}
           
function submit(inputValues, errorCode) {
  var ui = SpreadsheetApp.getUi();
  
  // 에러가 발생했는지 확인
  if (errorCode != '') {
    // alert 출력
    var msg = (errorCode in ErrorCodeMessages) ? ErrorCodeMessages[errorCode] : ErrorCodeMessages[undefined];
    ui.alert(msg);
    return;
  }
  
  // 변수들이 모두 입력되었는지 확인
  if (inputValues.profName === '' ||
      inputValues.classCode === '' ||
      inputValues.className === '' ||
      inputValues.numPeople === '') {
      ui.alert('누락된 정보가 있습니다. 모두 입력해주세요.');
  }
  
  // 백엔드에다가 여기 예약한다고 보내기
  
  // 최종 확인
  var msg = '예약 대상이 [' + inputValues.classRoom + ' / ' + inputValues.time + '] 이 맞습니까?';
  var response = ui.alert('예약 확인', msg, ui.ButtonSet.YES_NO);
  if (response != ui.Button.YES) {
    ui.alert('예약이 취소되었습니다.');
    return;
  }
  
  // 셀 병합
  var sheet = SpreadsheetApp.getActiveSheet();
  var range = getRange()[0];
  var sheetRange = sheet.getRange(range.row, range.col, range.height, range.width)
  sheetRange.merge();
  
  // 셀에 값 주기
  var cellValue = inputValues.className + '\n' + inputValues.profName;
  sheetRange.setValue(cellValue);
  
  ui.alert('예약됐다 치자.');
  return;
}