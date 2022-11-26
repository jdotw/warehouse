using Stock.Models;

namespace Stock.Services;

public static class CategoryService
{
  static List<Category> Categories { get; }
  static CategoryService()
  {
    Categories = new List<Category>
        {
            new Category { Id = Guid.NewGuid().ToString(), Name = "Diary" },
            new Category { Id = Guid.NewGuid().ToString(), Name = "Vegetables" }
        };
  }

  public static List<Category> GetAll() => Categories;

  public static Category? Get(string id) => Categories.FirstOrDefault(p => p.Id == id);

  public static void Add(Category Category)
  {
    Category.Id = Guid.NewGuid().ToString();
    Categories.Add(Category);
  }

  public static void Delete(string id)
  {
    var Category = Get(id);
    if (Category is null)
      return;
    Categories.Remove(Category);
  }

  public static void Update(Category Category)
  {
    var index = Categories.FindIndex(p => p.Id == Category.Id);
    if (index == -1)
      return;
    Categories[index] = Category;
  }
}