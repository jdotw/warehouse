namespace Stock.Models;

public class Category
{
  public string Id { get; set; }
  public string? Name { get; set; }

  public Category()
  {
    Id = Guid.NewGuid().ToString();
  }
}