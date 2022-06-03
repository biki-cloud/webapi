(function(document) {
    "use strict";

    var LightTableFilter = (function(Arr) {
        var _input;

        function _onInputEvent(e) {
            _input = e.target;
            var tables = document.getElementsByClassName(
                _input.getAttribute("data-table")
            );
            Arr.forEach.call(tables, function(table) {
                Arr.forEach.call(table.tBodies, function(tbody) {
                    Arr.forEach.call(tbody.rows, _filter);
                });
            });
        }

        function _filter(row) {
            // タイトルの行は残す
            if (row.getElementsByClassName("table-title").length > 0) {
                return;
            }

            // 1レコード中の文字列（submitボタンのバリューは除く)
            var textOfRecord = row.textContent.toLowerCase();

            // 検索文字列
            var seachingLetters = _input.value.toLowerCase();

            // 表の１行の中からinputタグを全て取得し、その中のsubmitボタンのvalueを取得する
            var submitValue = "";
            var inputTags = row.getElementsByTagName("input");
            if (inputTags.length > 0) {
                for (var i = 0; i < inputTags.length; i++) {
                    if (inputTags[i].type === "submit") {
                        submitValue = inputTags[i].value.toLowerCase();
                    }
                }
            }

            // 1レコードの文字列とsubmitボタンのバリューを検索し、
            // どちらにも検索文字列が存在しなかったら、そのレコードは見えなくする。
            if (
                textOfRecord.indexOf(seachingLetters) === -1 &&
                submitValue.indexOf(seachingLetters) === -1
            ) {
                row.style.display = "none";
            } else {
                row.style.display = "table-row";
            }
        }

        return {
            init: function() {
                var inputs = document.getElementsByClassName("light-table-filter");
                Arr.forEach.call(inputs, function(input) {
                    input.oninput = _onInputEvent;
                });
            },
        };
    })(Array.prototype);

    document.addEventListener("readystatechange", function() {
        if (document.readyState === "complete") {
            LightTableFilter.init();
        }
    });
})(document);