unit LogsDeleting;

{$mode ObjFPC}{$H+}

interface

procedure DeleteLogs(LogsDir : string);

implementation

procedure DeleteLogs(LogsDir : string);
begin
  Writeln (LogsDir, ' is a directory');
end;

end.

