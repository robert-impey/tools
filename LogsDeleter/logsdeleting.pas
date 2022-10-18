unit LogsDeleting;

{$mode ObjFPC}{$H+}

interface

uses
    Classes, FileUtil, SysUtils;

procedure DeleteLogs(LogsDir : string);
procedure DeleteLogsFromSubDir(SubDirPath: string);

implementation

procedure DeleteLogs(LogsDir : string);
var
   SearchPath, SubDirsPath : string;
   Info : TSearchRec;
   Count : longint;
begin
  Writeln(Format('Searching for subdirectories of %s', [LogsDir]));

  SearchPath := Concat(LogsDir, '/*');
  DoDirSeparators(SearchPath);
  Count:=0;
  If FindFirst (SearchPath,faAnyFile,Info)=0 then
    begin
    Repeat
      With Info do
        begin
          If ((Attr and faDirectory) = faDirectory) and (Name<>'.') and (Name<>'..') then
          begin
               SubDirsPath := Concat(LogsDir, '/', Name);
               DoDirSeparators(SubDirsPath);
               DeleteLogsFromSubDir(SubDirsPath);
               Inc(Count);
          end;
        end;
    Until FindNext(Info)<>0;
    FindClose(Info);
    end;
  Writeln ('Finished search. Found ',Count,' matches');
end;

procedure DeleteLogsFromSubDir(SubDirPath: string);
begin
     Writeln(Format('Searching for log files in %s', [SubDirPath]));
end;

end.

