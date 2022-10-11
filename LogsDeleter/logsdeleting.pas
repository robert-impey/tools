unit LogsDeleting;

{$mode ObjFPC}{$H+}

interface

procedure DeleteLogs(LogsDir : string);

implementation

uses
  Classes, FileUtil, SysUtils;

procedure DeleteLogs(LogsDir : string);
var
   SearchPath : string;
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
            Writeln(Name);
            Inc(Count);
          end;
        end;
    Until FindNext(info)<>0;
    FindClose(Info);
    end;
  Writeln ('Finished search. Found ',Count,' matches');
end;

end.

