using Microsoft.EntityFrameworkCore;
using Stock.Models;

namespace Stock.PostgreSQL
{
  public class StockContext : DbContext
  {
    public DbSet<Category> Categories { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        => optionsBuilder.UseNpgsql("Host=10.1.1.2;Database=stock_dotnet;Username=local;Password=local");
  }
}

