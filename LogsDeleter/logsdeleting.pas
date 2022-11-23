unit LogsDeleting;

{$mode ObjFPC}{$H+}

interface

uses
  Classes, FileUtil, SysUtils, DateUtils;

procedure DeleteLogs(LogsDir, Tool: string);
procedure DeleteLogsFromSubDir(SubDirPath: string);

implementation

procedure DeleteLogs(LogsDir, Tool: string);
var
  SearchPath, SubDirsPath: string;
  Info: TSearchRec;
begin
  if Tool = '' then
  begin
    Writeln(Format('Searching for subdirectories of %s', [LogsDir]));

    SearchPath := Concat(LogsDir, '/*');
    DoDirSeparators(SearchPath);
    if FindFirst(SearchPath, faAnyFile, Info) = 0 then
    begin
      repeat
        with Info do
        begin
          if ((Attr and faDirectory) = faDirectory) and (Name <> '.') and (Name <> '..') then
          begin
            SubDirsPath := Concat(LogsDir, '/', Name);
            DoDirSeparators(SubDirsPath);
            DeleteLogsFromSubDir(SubDirsPath);
          end;
        end;
      until FindNext(Info) <> 0;
      FindClose(Info);
    end ;
    end
    else
    begin
      SubDirsPath := Concat(LogsDir, '/', Tool);
      DoDirSeparators(SubDirsPath);
      DeleteLogsFromSubDir(SubDirsPath);
    end;
end;

procedure DeleteLogsFromSubDir(SubDirPath: string);
var
  Threshold, LogFileDate: TDateTime;
  LogFiles: TStringList;
  I: integer;
  FA: longint;
begin
  Writeln(Format('Searching for log files in %s', [SubDirPath]));

  Threshold := IncMonth(Date, -1);

  LogFiles := FindAllFiles(SubDirPath, '*.log;*.err', False);
  try
    Writeln(Format('Found %d Log files', [LogFiles.Count]));
    for I := 0 to pred(LogFiles.Count) do
    begin
      FA := FileAge(LogFiles[I]);
      if FA <> -1 then
      begin
        LogFileDate := FileDateToDateTime(FA);
        if DateOf(LogFileDate) <= DateOf(Threshold) then
        begin
          Writeln(Format('%s is old - deleting...', [LogFiles[I]]));
          DeleteFile(LogFiles[I]);
        end;
      end;
    end;
  finally
    LogFiles.Free;
  end;
end;

end.
