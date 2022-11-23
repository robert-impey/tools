program LogsDeleter;

{$mode objfpc}{$H+}

uses
 {$IFDEF UNIX}
  cthreads,
      {$ENDIF}
  Classes,
  SysUtils,
  CustApp,
  LogsDeleting;

type

  { TLogsDeleter }

  TLogsDeleter = class(TCustomApplication)
  protected
    procedure DoRun; override;
  public
    constructor Create(TheOwner: TComponent); override;
    destructor Destroy; override;
    procedure WriteHelp; virtual;
  end;

  { TLogsDeleter }

  procedure TLogsDeleter.DoRun;
  var
    ErrorMsg, LogsDir, Tool: string;
    LogsDirAttr: longint;
  begin
    // quick check parameters
    ErrorMsg := CheckOptions('ht', 'help tool');
    if ErrorMsg <> '' then
    begin
      ShowException(Exception.Create(ErrorMsg));
      Terminate;
      Exit;
    end;

    // parse parameters
    if HasOption('h', 'help') then
    begin
      WriteHelp;
      Terminate;
      Exit;
    end;

    if HasOption('t', 'tool') then
      Tool := GetOptionValue('t', 'tool')
    else
      Tool := '';

    LogsDir := Concat(GetUserDir, 'logs');
    DoDirSeparators(LogsDir);
    LogsDirAttr := FileGetAttr(LogsDir);

    if LogsDirAttr < 0 then
    begin
      WriteLn('Unable to read the file attributes of ', LogsDir);
      WriteLn('Does it exist?');
    end
    else
    if (LogsDirAttr and faDirectory) <> 0 then
      LogsDeleting.DeleteLogs(LogsDir, Tool);

    // stop program loop
    Terminate;
  end;

  constructor TLogsDeleter.Create(TheOwner: TComponent);
  begin
    inherited Create(TheOwner);
    StopOnException := True;
  end;

  destructor TLogsDeleter.Destroy;
  begin
    inherited Destroy;
  end;

  procedure TLogsDeleter.WriteHelp;
  begin
    { add your help code here }
    writeln('Usage: ', ExeName, ' -h');
  end;

var
  Application: TLogsDeleter;
begin
  Application := TLogsDeleter.Create(nil);
  Application.Title := 'Logs Deleter';
  Application.Run;
  Application.Free;
end.
