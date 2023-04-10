$("#loginToVendorBtn").click(postLogin);
$("#forceCategoryRefresh").click(forceCategoriesRefresh);
$("#populateProductUrlsBtn").click(populateProductUrls);

var loadedCategories;
$(document).ready(function () {
    var savedUser = window.localStorage.getItem('vendor_username');
    var savedPass = window.localStorage.getItem('vendor_password');
    if (savedUser && savedPass) {
        $("#argyrosEmailInput").val(savedUser);
        $("#argyrosPasswordInput").val(savedPass);
        postLogin();
    }
});

function postLogin() {
    btnToLoader($("#loginToVendorBtn"));

    var email = $("#argyrosEmailInput").val();
    var pass = $("#argyrosPasswordInput").val();

    $.post("/auth/vendor_login", JSON.stringify({ email: email, password: pass }), function (data) {
        $("#loginCard .card-header").addClass("bg-info bg-gradient");
        $("#loginCard .card-header h4").append(" ✔️");
        $("#loginCard .card-body").slideUp();

        window.localStorage.setItem('vendor_username', email);
        window.localStorage.setItem('vendor_password', pass);

        loadCategories()
    });
}


function loadCategories() {
    $("#getCategoriesCard").slideDown();

    var cachedCategories = window.localStorage.getItem('cached_categories');
    if (cachedCategories) {
        var categories = JSON.parse(cachedCategories);
        displayLoadedCategories(categories, true);
        return;
    }

    refreshCategories();
}


function displayLoadedCategories(categories, retrievedFromCache) {
    loadedCategories = categories.category_map;

    $("#subcategoriesLoaderContainer").slideUp();
    $("#getCategoriesLoadedContainer").slideDown();


    var html = "";
    $.each(loadedCategories, function (categoryName, subcategories) {
        html += `   
            <div class="card" style="height: 100%;">
            <div class="card-header">` + categoryName + `</div>
            <ul class="list-group list-group-flush">`;

        $.each(subcategories, function (subcategoryName, subcategoryUrl) {
            html += `<li class="list-group-item subcategory-card" cat-name="` + categoryName + `" subcat-name="` + subcategoryName + `" subcat-url="` + subcategoryUrl + `">` + subcategoryName + `</li>`;
        })
        html += `</ul></div>`;
    });
    $("#categoriesDisplayCardsContainer").html(html);

    if (retrievedFromCache) {
        $("#forceRefreshSubcategoriesContainer").slideDown();
    }
}

function forceCategoriesRefresh() {
    $("#subcategoriesLoaderContainer").slideDown();
    $("#getCategoriesLoadedContainer").slideUp();
    $("#forceRefreshSubcategoriesContainer").slideUp();
    refreshCategories();
}

function refreshCategories() {
    $.get("/api/all_categories", function (data) {
        var cachedCategoriesJson = JSON.stringify(data)
        window.localStorage.setItem('cached_categories', cachedCategoriesJson);
        displayLoadedCategories(data, false);
    });
};


function populateProductUrls() {

    loadedCatalogProducts = {};

    var totalProducts = 0;
    $(".subcategory-card").each(function (index) {
        var category = $(this).attr("cat-name");
        var subcategory = $(this).attr("subcat-name");
        var subcategoryURL = $(this).attr("subcat-url");
        var onlyLast = $("#categoryDownloadLastPageCheck").is(":checked");

        var subCatContainer = $(this);
        subCatContainer.append(loaderHtml);

        $.get("/api/subcategory_products?subcat_url=" + encodeURIComponent(subcategoryURL) + "&last_page=" + onlyLast, function (data) {
            // FOR debug
            data.products_urls = data.products_urls;

            subCatContainer.find(".spinner-border").remove();
            subCatContainer.addClass("bg-success bg-gradient");
            subCatContainer.append(" (" + data.products_urls.length + ")");
            totalProducts += data.products_urls.length;

            if (!(category in loadedCatalogProducts)) {
                loadedCatalogProducts[category] = {};
            }
            loadedCatalogProducts[category][subcategory] = data.products_urls;

            var remainingLoaders = $(".subcategory-card .spinner-border").length;
            if (remainingLoaders == 0) {
                $("#getCategoriesCard .card-header").addClass("bg-info bg-gradient");
                $("#getCategoriesCard .card-header h4").append(" ✔️");
                $("#getCategoriesCard .card-body").slideUp();

                showDownloadImages(totalProducts);
            }
        });
    });
}

var onlyNewProductsCount = 0;
function showDownloadImages(productsCount) {
    $(".fetchedProductsCount").html(productsCount);

    var lastFetchTimestamp = window.localStorage.getItem('prev_downloaded_timestamp');
    if (lastFetchTimestamp) {
        $("#fetchedProductsLastExecution").html(lastFetchTimestamp);
        loadedCatalogOnlyNewProducts = calculateOnlyNewProducts();
        $(".fetchedProductsCountUnique").html(onlyNewProductsCount);

        if(onlyNewProductsCount == 0){
            $("#onlyNewImagesContainer").html("No new items since last download on <strong>"+lastFetchTimestamp+"</strong>")
        }

    } else {
        $("#onlyNewImagesContainer").html("No previous execution recorded. So no partial download available yet.")
    }



    $("#downloadImagesContainer").slideDown();
}


var loaderHtml = `<div class="spinner-border text-info" role="status">
<span class="visually-hidden">Loading...</span>
</div>`;

function btnToLoader(button) {
    $(button).html(loaderHtml);
    $(button).addClass("disabled");
}

function updateProcessedCategoriesRecord() {
    var prevCatSavedStr = window.localStorage.getItem('prev_downloaded_images_record');
    if (!prevCatSavedStr) {
        window.localStorage.setItem('prev_downloaded_images_record', JSON.stringify(loadedCatalogProducts));
        return;
    }

    var prevCatSaved = JSON.parse(prevCatSavedStr);

    // Merge loadedCatalogProducts with prevCatSaved
    $.each(loadedCatalogProducts, function (categoryName, subcategories) {
        if (!(categoryName in prevCatSaved)) {
            prevCatSaved[categoryName] = {};
        }
        $.each(subcategories, function (subcategoryName, subcategoryUrl) {
            if (!(subcategoryName in prevCatSaved[categoryName])) {
                prevCatSaved[categoryName][subcategoryName] = loadedCatalogProducts[categoryName][subcategoryName];
            } else {
                var newUrlsArr = prevCatSaved[categoryName][subcategoryName].concat(loadedCatalogProducts[categoryName][subcategoryName]);
                prevCatSaved[categoryName][subcategoryName] = [...new Set(newUrlsArr)];
            }
        });
    })

    window.localStorage.setItem('prev_downloaded_images_record', JSON.stringify(prevCatSaved));

    var currentdate = new Date();
    var datetime = currentdate.getDate() + "/"
        + (currentdate.getMonth() + 1) + "/"
        + currentdate.getFullYear() + " "
        + currentdate.getHours() + ":"
        + currentdate.getMinutes() + ":"
        + currentdate.getSeconds();
    window.localStorage.setItem('prev_downloaded_timestamp', datetime);

}


function calculateOnlyNewProducts() {
    var prevCatSavedStr = window.localStorage.getItem('prev_downloaded_images_record');
    if (!prevCatSavedStr) {
        // If no prev, all is new
        return loadedCatalogProducts;
    }

    var prevCatSaved = JSON.parse(prevCatSavedStr);
    var onlyNewItemsCat = {}
    var totalNewProductsCount = 0;
    // Merge loadedCatalogProducts with prevCatSaved
    $.each(loadedCatalogProducts, function (categoryName, subcategories) {
        $.each(subcategories, function (subcategoryName, subcategoryUrl) {
            var newProductUrls = [];

            subcategoryUrl.forEach(productUrl => {
                if (!(categoryName in prevCatSaved) ||
                    !(subcategoryName in prevCatSaved[categoryName]) ||
                    $.inArray(productUrl, prevCatSaved[categoryName][subcategoryName]) === -1)  {
                    newProductUrls.push(productUrl);
                }
            });

            if (newProductUrls.length > 0) {
                if (!(categoryName in onlyNewItemsCat)) {
                    onlyNewItemsCat[categoryName] = {}
                }
                
                onlyNewItemsCat[categoryName][subcategoryName] = newProductUrls
                totalNewProductsCount += newProductUrls.length;
            }
        });
    })

    onlyNewProductsCount = totalNewProductsCount;

    return onlyNewItemsCat;
}