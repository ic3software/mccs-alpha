const getCounties = country => {
    const mapping = {
        England: ["Bedfordshire","Berkshire","Bristol","Buckinghamshire","Cambridgeshire","Cheshire","City of London","Cornwall","County Durham","Cumbria","Derbyshire","Devon","Dorset","East Riding of Yorkshire","East Sussex","Essex","Gloucestershire","Greater London","Greater Manchester","Hampshire","Herefordshire","Hertfordshire","Humberside","Isle of Wight","Isles of Scilly","Kent","Lancashire","Leicestershire","Lincolnshire","Merseyside","Middlesex","Norfolk","North Somerset","North Yorkshire","Northamptonshire","Northumberland","Nottinghamshire","Oxfordshire","Rutland","Shropshire","Somerset","South Gloucestershire","South Yorkshire","Staffordshire","Suffolk","Surrey","Tyne & Wear","Warwickshire","West Midlands","West Sussex","West Yorkshire","Wiltshire","Worcestershire"],
        Scotland: ["Aberdeenshire","Angus","Argyll & Bute","Ayrshire","Banffshire","Berwickshire","Borders","Caithness","Clackmannanshire","Dumfries & Galloway","Dunbartonshire","East Ayrshire","East Dunbartonshire","East Lothian","East Renfrewshire","Fife","Highland","Inverclyde","Kincardineshire","Lanarkshire","Midlothian","Moray","North Ayrshire","North Lanarkshire","Orkney","Perth & Kinross","Renfrewshire","Shetland","South Ayrshire","South Lanarkshire","Stirlingshire","West Dunbartonshire","West Lothian","Western Isles"],
        Wales: ["Blaenau Gwent","Bridgend","Caerphilly","Cardiff","Carmarthenshire","Ceredigion","Conwy","Denbighshire","Flintshire","Gwynedd","Isle of Anglesey","Merthyr Tydfil","Monmouthshire","Neath Port Talbot","Newport","Pembrokeshire","Powys","Rhondda Cynon Taff","Swansea","Torfaen","Vale of Glamorgan","Wrexham"],
        "Northern Ireland": ["Antrim","Armagh","Down","Fermanagh","Londonderry","Tyrone"]
    }
    return mapping[country]
}


$(function () {
    // ************ Check user membership ****************
    $.ajax({
        url: `/api/is-trading-member`,
        method: "GET",
        success: data => {
            if (data.IsMember) {
                $(".header-transfer-link").show()
                $(".header-history-link").show()
            }
        },
        error: function () {}
    });
    // ****************************

    // ************ Dropdown ****************
    $(".ui.dropdown").dropdown();

    $(".multiple.dropdown").dropdown({
        allowAdditions: true
    });

    $(".wants.dropdown, .offers.dropdown").dropdown({
        allowAdditions: true,
        ignoreCase: true,
        forceSelection: false,
        hideAdditions: false,
        apiSettings: {
            throttle: 700,
            url: '/api/tags/{query}',
            beforeSend: function (settings) {
                if (settings.urlData.query === "") {
                    return false;
                }
                return settings
            }
        },
    });
    // ****************************

    // ************ Advanced Search ****************
    if (localStorage.getItem('advanced-search') === "true") {
        $("#advanced-search").css("display", "")
        $("#search-toggle").text("Basic Search")
    }
    // ****************************

    $(document).on("keydown", ":input:not(textarea)", function (event) {
        return event.key != "Enter";
    });

    // Delete Account
    $(".action-delete").click(function () {
        const bID = $(this).attr("business-id")
        $(".ui.basic.modal").modal("show");
        $("#selectedId").val(bID)
    });

    $("#confirm-delete").click(function () {
        const bID = $("#selectedId").val()
        $.ajax({
            url: `/admin/api/businesses/${bID}`,
            method: "DELETE",
            success: function () {
                $(`#${bID}`).remove()
            },
            error: function (xhr) {
                alert("An error occurred: " + xhr.status + " " + xhr.statusText);
            }
        });
    });

    // Create Tag
    $("#submit-tag").click(function () {
        const tagName = $("#new-tag").val()
        $.ajax({
            url: "/admin/api/user-tags",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                name: tagName
            }),
            success: function () {
                showSuccessMessage("Tag created.")
            },
            error: function (xhr) {
                showErrorMessage(xhr.responseText);
            }
        });
    });

    $("#confirm-delete-tag").click(function () {
        const tID = $("#selectedId").val()
        $.ajax({
            url: `/admin/api/user-tags/${tID}`,
            method: "DELETE",
            success: function () {
                $(`#${tID}`).remove()
                showSuccessMessage("Tag has been removed.")
            },
            error: function (xhr) {
                showErrorMessage("An error occurred: " + xhr.status + " " + xhr.statusText);
            }
        });
    });

    // Create Admin Tag
    $("#submit-admin-tag").click(function () {
        let tagName = $("#new-admin-tag").val()
        tagName = tagName.replace(/[0-9]|(&quot;)|([^a-zA-Z ]+)/g, "")
        $.ajax({
            url: "/admin/api/admin-tags",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                name: tagName
            }),
            success: function () {
                showSuccessMessage("Admin tag created.")
            },
            error: function (xhr) {
                showErrorMessage(xhr.responseText);
            }
        });
    });

    // ************ Update Tag ****************
    $(".action-update-tag").click(function () {
        const tID = $(this).attr("tag-id")
        $(".ui.basic.update.modal").modal("show");
        $("#selectedId").val(tID)
    });

    $("#confirm-update-tag").click(function () {
        const id = $("#selectedId").val()
        const name = $(`input[tag-id=${id}]`).val();
        $.ajax({
            url: `/admin/api/user-tags/${id}`,
            method: "PUT",
            contentType: "application/json",
            data: JSON.stringify({
                id: id,
                name: name
            }),
            success: function () {
                showSuccessMessage("Tag updated.")
                $(`tr[id=${id}] td:first`).html(name);
            },
            error: function (xhr) {
                showErrorMessage(xhr.responseText);
            }
        });
    });
    // ****************************

    // ************ Update Admin Tag ****************
    $(".action-update-admin-tag").click(function () {
        const tID = $(this).attr("admin-tag-id")
        $(".ui.basic.update.modal").modal("show");
        $("#selectedId").val(tID)
    });

    $("#confirm-update-admin-tag").click(function () {
        const id = $("#selectedId").val()
        let name = $(`input[admin-tag-id=${id}]`).val();
        name = name.replace(/[0-9]|(&quot;)|([^a-zA-Z ]+)/g, "")
        $.ajax({
            url: `/admin/api/admin-tags/${id}`,
            method: "PUT",
            contentType: "application/json",
            data: JSON.stringify({
                id: id,
                name: name
            }),
            success: function () {
                showSuccessMessage("Tag updated.")
                $(`tr[id=${id}] td:first`).html(name);
                $(`input[admin-tag-id=${id}]`).val("");
            },
            error: function (xhr) {
                showErrorMessage(xhr.responseText);
            }
        });
    });
    // ****************************

    // ************ Delete Tag ****************
    $(".action-delete-tag").click(function () {
        const tID = $(this).attr("tag-id")
        $(".ui.basic.delete.modal").modal("show");
        $("#selectedId").val(tID)
    });
    // ****************************

    // ************  Delete Admin Tag ****************
    $(".action-delete-admin-tag").click(function () {
        const tID = $(this).attr("admin-tag-id")
        $(".ui.basic.delete.modal").modal("show");
        $("#selectedId").val(tID)
    });

    $("#confirm-delete-admin-tag").click(function () {
        const tID = $("#selectedId").val()
        $.ajax({
            url: `/admin/api/admin-tags/${tID}`,
            method: "DELETE",
            success: function () {
                $(`#${tID}`).remove()
                showSuccessMessage("Tag has been removed.")
            },
            error: function (xhr) {
                showErrorMessage("An error occurred: " + xhr.status + " " + xhr.statusText);
            }
        });
    });
    // ****************************
});

const toggleMenu = () => $(".ui.sidebar").sidebar("toggle");

// ToggleFavoriteBusinesses
function handleClickFavorite(bID) {
    const elem = $(`.favorite-${bID}`)
    if (elem.hasClass("outline")) {
        $.ajax({
            url: `/api/users/addToFavoriteBusinesses`,
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                id: bID
            }),
            success: function() {
                elem.removeClass("outline")
                elem.addClass("red")
            },
            error: function(xhr) {
                showErrorMessage("An error occurred: " + xhr.status + " " + xhr.statusText);
            }
        });
    } else {
        $.ajax({
            url: `/api/users/removeFromFavoriteBusinesses`,
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                id: bID
            }),
            success: function() {
                elem.addClass("outline")
                elem.removeClass("red")
            },
            error: function(xhr) {
                showErrorMessage("An error occurred: " + xhr.status + " " + xhr.statusText);
            }
        });
    }
}

const showSuccessMessage = message => {
    $(".ui.success").addClass("hidden")
    $(".ui.error").addClass("hidden")
    $(".ajax-success").removeClass("hidden")
    $(".ajax-success .header").text(message)
}

const showErrorMessage = message => {
    $(".ui.success").addClass("hidden")
    $(".ui.error").addClass("hidden")
    $(".ajax-error").removeClass("hidden")
    $(".ajax-error .header").text(message)
}

// Advanced Search Toggle
const handleClickAdvancedSearch = () => {
    if ($("#advanced-search").css("display") === "none") {
        localStorage.setItem("advanced-search", true)
        $("#search-toggle").text("Basic Search")
    } else {
        localStorage.setItem("advanced-search", false)
        $("#search-toggle").text("Advanced Search")
    }
    $("#advanced-search").toggle(350, () => {})
};

function debounce(func, wait, immediate) {
	var timeout;
	return function() {
		var context = this, args = arguments;
		var later = function() {
			timeout = null;
			if (!immediate) func.apply(context, args);
		};
		var callNow = immediate && !timeout;
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
		if (callNow) func.apply(context, args);
	};
};

const formatTime = dateString => {
    const appendLeadingZeroes = n => {
        if (n <= 9) {
          return "0" + n;
        }
        return n
    }

    const date = new Date(dateString)
    let hh = appendLeadingZeroes(date.getUTCHours());
    let mm = appendLeadingZeroes(date.getUTCMinutes());
    let ss = appendLeadingZeroes(date.getSeconds());

    return `${date.getFullYear()}-${appendLeadingZeroes(date.getMonth() + 1)}-${appendLeadingZeroes(date.getDate())} ${hh}:${mm}:${ss} UTC`
}
