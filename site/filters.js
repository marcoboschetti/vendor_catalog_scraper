var categories;

var fullSubcategoriesList;
$(document).ready(function () {
    loadCategories();
});

function loadCategories() {
    $.get("/api/admin/categories", function (data) {
        categories = data.categories;
        categories.sort(function SortByName(a, b) {
            var aName = a.name.toLowerCase();
            var bName = b.name.toLowerCase();
            return ((aName < bName) ? -1 : ((aName > bName) ? 1 : 0));
        });

        var categoriesListHtml = "";
        categories.forEach(category => {
            categoriesListHtml += `<li class="category-filter-item" cat-id=`+category.id+`><a href="#">` + category.name + `</a></li>`;
        });
        $("#shopFiltersCategories").html(categoriesListHtml);
        $(".category-filter-item").first().addClass("selected");
        $(".category-filter-item").click(toggleCategoryFilter)
        populateSubcategoriesFilters(categories[0].id);
    });
}

function toggleCategoryFilter(){
    var item = $(this);
    $(".category-filter-item").removeClass("selected");
    item.addClass("selected");
    var categoryID = item.attr("cat-id")
    populateSubcategoriesFilters(categoryID); 
}

function populateSubcategoriesFilters(categoryID){
    var category;
    categories.forEach(c => {
        if(c.id == categoryID){
            category = c;
        }
    });

    var subcategoriesListHtml = "";
    category.subcategories.forEach(subcategory => {
        subcategoriesListHtml += `<li class="subcategory-filter-item" cat-id=`+category.id+` subcat-id=`+subcategory.id+`><a href="#">` + subcategory.name + `</a></li>`;
    });
    $("#shopFiltersSubcategories").html(subcategoriesListHtml);
    $(".subcategory-filter-item").first().addClass("selected");
    $(".subcategory-filter-item").click(toggleSubcategoryFilter);
    loadItemsFromFilters();
}

function toggleSubcategoryFilter(){
    var item = $(this);
    $(".subcategory-filter-item").removeClass("selected");
    item.addClass("selected");
    loadItemsFromFilters();
}
