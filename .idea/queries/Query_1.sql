SELECT
    soh.SalesOrderNumber AS OrderNumber,
    COUNT(DISTINCT pc.ProductCategoryID) AS CategoriesCount,
    COUNT(DISTINCT ps.ProductSubcategoryID) AS SubcategoriesCount
FROM
    Sales.SalesOrderHeader soh
        INNER JOIN Sales.SalesOrderDetail sod ON soh.SalesOrderID = sod.SalesOrderID
        INNER JOIN Production.Product p ON sod.ProductID = p.ProductID
        INNER JOIN Production.ProductSubcategory ps ON p.ProductSubcategoryID = ps.ProductSubcategoryID
        INNER JOIN Production.ProductCategory pc ON ps.ProductCategoryID = pc.ProductCategoryID
GROUP BY
    soh.SalesOrderID,
    soh.SalesOrderNumber
ORDER BY CategoriesCount DESC;

WITH
    ProductSales AS (
        SELECT
            p.ProductID,
            p.Name AS ProductName,
            COUNT(*) AS TotalSales,
            COUNT(DISTINCT soh.CustomerID) AS ProductCustomers
        FROM
            Production.Product p
                JOIN Sales.SalesOrderDetail sod ON p.ProductID = sod.ProductID
                JOIN Sales.SalesOrderHeader soh ON sod.SalesOrderID = soh.SalesOrderID
        GROUP BY
            p.ProductID, p.Name
    ),
    CategoryCustomers AS (
        SELECT
            pc.ProductCategoryID,
            COUNT(DISTINCT soh.CustomerID) AS CategoryCustomers
        FROM
            Production.ProductCategory pc
                JOIN Production.ProductSubcategory ps ON pc.ProductCategoryID = ps.ProductCategoryID
                JOIN Production.Product p ON ps.ProductSubcategoryID = p.ProductSubcategoryID
                JOIN Sales.SalesOrderDetail sod ON p.ProductID = sod.ProductID
                JOIN Sales.SalesOrderHeader soh ON sod.SalesOrderID = soh.SalesOrderID
        GROUP BY
            pc.ProductCategoryID
    )
SELECT
    ps.ProductID,
    ps.ProductName,
    ps.TotalSales,
    ps.ProductCustomers,
    CAST(ps.ProductCustomers AS DECIMAL) / cc.CategoryCustomers AS CustomerRatio
FROM
    ProductSales ps
        JOIN Production.ProductSubcategory psc ON ps.ProductID = psc.ProductSubcategoryID
        JOIN Production.ProductCategory pc ON psc.ProductCategoryID = pc.ProductCategoryID
        JOIN CategoryCustomers cc ON pc.ProductCategoryID = cc.ProductCategoryID;