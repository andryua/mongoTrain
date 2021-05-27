$(function() {
    $("#tree").fancytree({
        checkbox: true,
        selectMode: 3,
        extensions: ["filter"],
        quicksearch: true,
        source: {
            url: "/admjson",
        },
        filter: {
            autoApply: true, // Re-apply last filter if lazy data is loaded
            autoExpand: false, // Expand all branches that contain matches while filtered
            counter: true, // Show a badge with number of matching child nodes near parent icons
            fuzzy: false, // Match single characters in order, e.g. 'fb' will match 'FooBar'
            hideExpandedCounter: true, // Hide counter badge if parent is expanded
            hideExpanders: false, // Hide expanders if all child nodes are hidden by filter
            highlight: true, // Highlight matches by wrapping inside <mark> tags
            leavesOnly: false, // Match end nodes only
            nodata: true, // Display a 'no data' status node if result is empty
            mode: "hide" // Grayout unmatched nodes (pass "hide" to remove unmatched node instead)
        },

        lazyLoad: function (event, ctx) {
            ctx.result = {url: "filter_result.json"};
        },
        loadChildren: function (event, ctx) {
            ctx.node.fixSelection3AfterClick();
        },
        activate: function (event, data) {
            if (data.node.title && !data.node.folder) {
                $("#tit").show();
                $("#echoName").text(data.node.title)
                $("#tit1").show();
                $("#tit2").show();
                $("#tit3").show();
                $("#echoSupport").text(data.node.data.support);
                $("#echoInfo").text(data.node.data.description);
                $("#echoValues").text(JSON.stringify(data.node.data.values));
            } else {
                $("#tit").hide();
                $("#tit1").hide();
                $("#tit2").hide();
                $("#tit3").hide();
                $("#echoName").text("")
                $("#echoSupport").text("");
                $("#echoInfo").text("");
            }
        },
        deactivate: function () {
            $("#tit1").hide();
            $("#tit2").hide();
            $("#tit3").hide();
            $("#echoSupport").text("");
            $("#echoInfo").text("");
            $("#echoName").text("");
        },

        select: function (event, data) {
            var selKeys = $.map(data.tree.getSelectedNodes(), function (node) {
                return node.data.id;
            });
            //$("#ids").text(selKeys.join(","));
            //$("#ids").val(selKeys.join(","));
            //$("#admtmpid").val(function() {
            $("#ids").val(function () {
                var emphasis = selKeys.join(",");
                return emphasis;
            });
            //$("#echoSelectionRootKeys3").text(selRootKeys.join(", "));
            // $("#echoSelectionRoots3").text(selRootNodes.join(", "));
        },
        // The following options are only required, if we have more than one tree on one page:
        cookieId: "fancytree-Cb3",
        idPrefix: "fancytree-Cb3-"
    });

    var tree = $.ui.fancytree.getTree("#tree");

    //var tree = $("#tree").fancytree("getTree");

    $(".fancytree-container").addClass("fancytree-connectors");

    $("input[name=search]").on("keyup", function (e) {
        var n,
            tree = $.ui.fancytree.getTree(),
            args = "autoApply autoExpand fuzzy hideExpanders highlight leavesOnly nodata".split(" "),
            opts = {},
            filterFunc = tree.filterNodes,
            match = $(this).val();

        opts.mode = "hide";
        //opts.hideExpandedCounter = true;
        //opts.autoExpand = false;

        if (e && e.which === $.ui.keyCode.ESCAPE || $.trim(match) === "") {
            $("button#btnResetSearch").click();
            return;
        }
        n = filterFunc.call(tree, match, opts);
        $("button#btnResetSearch").attr("disabled", false);
    }).focus();

    $("button#btnResetSearch").click(function () {
        $("input[name=search]").val("");
        tree.clearFilter();
    }).attr("disabled", true);
});