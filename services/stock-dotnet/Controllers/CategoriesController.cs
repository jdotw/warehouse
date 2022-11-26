using Microsoft.AspNetCore.Mvc;
using Stock.Models;
using Stock.PostgreSQL;
using Stock.Services;
using Microsoft.EntityFrameworkCore;

namespace Stock.Controllers;

[ApiController]
[Route("[controller]")]
public class Categories : ControllerBase
{
  private readonly ILogger<Categories> _logger;
  private readonly StockContext _context;

  public Categories(StockContext context, ILogger<Categories> logger)
  {
    _logger = logger;
    _context = context;
  }

  [HttpGet]
  public ActionResult<List<Category>> GetAll()
  {
    var categories = _context.Categories.Select(c => new Category
    {
      Id = c.Id,
      Name = c.Name
    });
    return (categories == null) ? NotFound() : Ok(categories);
  }

  [HttpGet("{id}")]
  public ActionResult<Category> Get(string id)
  {
    // var category = CategoryService.Get(id);
    var category = _context.Categories.FirstOrDefault(c => c.Id == id);
    if (category == null)
      return NotFound();
    return category;
  }

  [HttpPost]
  public IActionResult Create(Category category)
  {
    _context.Categories.Add(category);
    _context.SaveChanges();
    return CreatedAtAction(nameof(Create), new { id = category.Id }, category);
  }

  [HttpPut("{id}")]
  public IActionResult Update(string id, Category category)
  {
    if (id != category.Id)
      return BadRequest();

    var existingCategory = _context.Categories.AsNoTracking().FirstOrDefault(c => c.Id == id);
    if (existingCategory is null)
      return NotFound();

    _context.Categories.Update(category);
    _context.SaveChanges();

    return NoContent();
  }
}

