ace.define("ace/mode/abap_highlight_rules",["require","exports","module","ace/lib/oop","ace/mode/text_highlight_rules"],function(E,e,t){"use strict";var N=E("../lib/oop"),T=E("./text_highlight_rules").TextHighlightRules,I=function(){var E=this.createKeywordMapper({"variable.language":"this",keyword:"ADD ALIAS ALIASES ASSERT ASSIGN ASSIGNING AT BACK CALL CASE CATCH CHECK CLASS CLEAR CLOSE CNT COLLECT COMMIT COMMUNICATION COMPUTE CONCATENATE CONDENSE CONSTANTS CONTINUE CONTROLS CONVERT CREATE CURRENCY DATA DEFINE DEFINITION DEFERRED DELETE DESCRIBE DETAIL DIVIDE DO ELSE ELSEIF ENDAT ENDCASE ENDCLASS ENDDO ENDEXEC ENDFORM ENDFUNCTION ENDIF ENDIFEND ENDINTERFACE ENDLOOP ENDMETHOD ENDMODULE ENDON ENDPROVIDE ENDSELECT ENDTRY ENDWHILE EVENT EVENTS EXEC EXIT EXPORT EXPORTING EXTRACT FETCH FIELDS FORM FORMAT FREE FROM FUNCTION GENERATE GET HIDE IF IMPORT IMPORTING INDEX INFOTYPES INITIALIZATION INTERFACE INTERFACES INPUT INSERT IMPLEMENTATION LEAVE LIKE LINE LOAD LOCAL LOOP MESSAGE METHOD METHODS MODIFY MODULE MOVE MULTIPLY ON OVERLAY OPTIONAL OTHERS PACK PARAMETERS PERFORM POSITION PROGRAM PROVIDE PUT RAISE RANGES READ RECEIVE RECEIVING REDEFINITION REFERENCE REFRESH REJECT REPLACE REPORT RESERVE RESTORE RETURNING ROLLBACK SCAN SCROLL SEARCH SELECT SET SHIFT SKIP SORT SORTED SPLIT STANDARD STATICS STEP STOP SUBMIT SUBTRACT SUM SUMMARY SUPPRESS TABLES TIMES TRANSFER TRANSLATE TRY TYPE TYPES UNASSIGN ULINE UNPACK UPDATE WHEN WHILE WINDOW WRITE OCCURS STRUCTURE OBJECT PROPERTY CASTING APPEND RAISING VALUE COLOR CHANGING EXCEPTION EXCEPTIONS DEFAULT CHECKBOX COMMENT ID NUMBER FOR TITLE OUTPUT WITH EXIT USING INTO WHERE GROUP BY HAVING ORDER BY SINGLE APPENDING CORRESPONDING FIELDS OF TABLE LEFT RIGHT OUTER INNER JOIN AS CLIENT SPECIFIED BYPASSING BUFFER UP TO ROWS CONNECTING EQ NE LT LE GT GE NOT AND OR XOR IN LIKE BETWEEN","constant.language":"TRUE FALSE NULL SPACE","support.type":"c n i p f d t x string xstring decfloat16 decfloat34","keyword.operator":"abs sign ceil floor trunc frac acos asin atan cos sin tan abapOperator cosh sinh tanh exp log log10 sqrt strlen xstrlen charlen numofchar dbmaxlen lines"},"text",!0," "),e="WITH\\W+(?:HEADER\\W+LINE|FRAME|KEY)|NO\\W+STANDARD\\W+PAGE\\W+HEADING|EXIT\\W+FROM\\W+STEP\\W+LOOP|BEGIN\\W+OF\\W+(?:BLOCK|LINE)|BEGIN\\W+OF|END\\W+OF\\W+(?:BLOCK|LINE)|END\\W+OF|NO\\W+INTERVALS|RESPECTING\\W+BLANKS|SEPARATED\\W+BY|USING\\W+(?:EDIT\\W+MASK)|WHERE\\W+(?:LINE)|RADIOBUTTON\\W+GROUP|REF\\W+TO|(?:PUBLIC|PRIVATE|PROTECTED)(?:\\W+SECTION)?|DELETING\\W+(?:TRAILING|LEADING)(?:ALL\\W+OCCURRENCES)|(?:FIRST|LAST)\\W+OCCURRENCE|INHERITING\\W+FROM|LINE-COUNT|ADD-CORRESPONDING|AUTHORITY-CHECK|BREAK-POINT|CLASS-DATA|CLASS-METHODS|CLASS-METHOD|DIVIDE-CORRESPONDING|EDITOR-CALL|END-OF-DEFINITION|END-OF-PAGE|END-OF-SELECTION|FIELD-GROUPS|FIELD-SYMBOLS|FUNCTION-POOL|MOVE-CORRESPONDING|MULTIPLY-CORRESPONDING|NEW-LINE|NEW-PAGE|NEW-SECTION|PRINT-CONTROL|RP-PROVIDE-FROM-LAST|SELECT-OPTIONS|SELECTION-SCREEN|START-OF-SELECTION|SUBTRACT-CORRESPONDING|SYNTAX-CHECK|SYNTAX-TRACE|TOP-OF-PAGE|TYPE-POOL|TYPE-POOLS|LINE-SIZE|LINE-COUNT|MESSAGE-ID|DISPLAY-MODE|READ(?:-ONLY)?|IS\\W+(?:NOT\\W+)?(?:ASSIGNED|BOUND|INITIAL|SUPPLIED)";this.$rules={start:[{token:"string",regex:"`",next:"string"},{token:"string",regex:"'",next:"qstring"},{token:"doc.comment",regex:/^\*.+/},{token:"comment",regex:/".+$/},{token:"invalid",regex:"\\.{2,}"},{token:"keyword.operator",regex:/\W[\-+\%=<>*]\W|\*\*|[~:,\.&$]|->*?|=>/},{token:"paren.lparen",regex:"[\\[({]"},{token:"paren.rparen",regex:"[\\])}]"},{token:"constant.numeric",regex:"[+-]?\\d+\\b"},{token:"variable.parameter",regex:/sy|pa?\d\d\d\d\|t\d\d\d\.|innnn/},{token:"keyword",regex:e},{token:"variable.parameter",regex:/\w+-\w+(?:-\w+)*/},{token:E,regex:"\\b\\w+\\b"},{caseInsensitive:!0}],qstring:[{token:"constant.language.escape",regex:"''"},{token:"string",regex:"'",next:"start"},{defaultToken:"string"}],string:[{token:"constant.language.escape",regex:"``"},{token:"string",regex:"`",next:"start"},{defaultToken:"string"}]}};N.inherits(I,T),e.AbapHighlightRules=I}),ace.define("ace/mode/folding/coffee",["require","exports","module","ace/lib/oop","ace/mode/folding/fold_mode","ace/range"],function(E,e,t){"use strict";var N=E("../../lib/oop"),T=E("./fold_mode").FoldMode,I=E("../../range").Range,n=e.FoldMode=function(){};N.inherits(n,T),function(){this.getFoldWidgetRange=function(E,e,t){var N=this.indentationBlock(E,t);if(N)return N;var T=/\S/,n=E.getLine(t),O=n.search(T);if(-1!=O&&"#"==n[O]){for(var R=n.length,r=E.getLength(),S=t,o=t;++t<r;){n=E.getLine(t);var i=n.search(T);if(-1!=i){if("#"!=n[i])break;o=t}}if(o>S){var a=E.getLine(o).length;return new I(S,R,o,a)}}},this.getFoldWidget=function(E,e,t){var N=E.getLine(t),T=N.search(/\S/),I=E.getLine(t+1),n=E.getLine(t-1),O=n.search(/\S/),R=I.search(/\S/);if(-1==T)return E.foldWidgets[t-1]=-1!=O&&R>O?"start":"","";if(-1==O){if(T==R&&"#"==N[T]&&"#"==I[T])return E.foldWidgets[t-1]="",E.foldWidgets[t+1]="","start"}else if(O==T&&"#"==N[T]&&"#"==n[T]&&-1==E.getLine(t-2).search(/\S/))return E.foldWidgets[t-1]="start",E.foldWidgets[t+1]="","";return E.foldWidgets[t-1]=-1!=O&&T>O?"start":"",R>T?"start":""}}.call(n.prototype)}),ace.define("ace/mode/abap",["require","exports","module","ace/mode/abap_highlight_rules","ace/mode/folding/coffee","ace/range","ace/mode/text","ace/lib/oop"],function(E,e,t){"use strict";function N(){this.HighlightRules=T,this.foldingRules=new I}var T=E("./abap_highlight_rules").AbapHighlightRules,I=E("./folding/coffee").FoldMode,n=E("../range").Range,O=E("./text").Mode,R=E("../lib/oop");R.inherits(N,O),function(){this.getNextLineIndent=function(E,e,t){var N=this.$getIndent(e);return N},this.toggleCommentLines=function(E,e,t,N){for(var T=new n(0,0,0,0),I=t;N>=I;++I){var O=e.getLine(I);hereComment.test(O)||(O=commentLine.test(O)?O.replace(commentLine,"$1"):O.replace(indentation,"$&#"),T.end.row=T.start.row=I,T.end.column=O.length+1,e.replace(T,O))}},this.$id="ace/mode/abap"}.call(N.prototype),e.Mode=N});