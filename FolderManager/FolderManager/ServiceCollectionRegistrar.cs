using Microsoft.Extensions.DependencyInjection;
using Spectre.Console.Cli;

namespace FolderManager;

public sealed class ServiceCollectionRegistrar : ITypeRegistrar
{
    private readonly IServiceCollection _services;

    public ServiceCollectionRegistrar(IServiceCollection services)
    {
        _services = services;
    }

    public ITypeResolver Build()
    {
        return new ServiceCollectionResolver(_services.BuildServiceProvider());
    }

    public void Register(Type service, Type implementation)
    {
        _services.AddSingleton(service, implementation);
    }

    public void RegisterInstance(Type service, object implementation)
    {
        _services.AddSingleton(service, implementation);
    }

    public void RegisterLazy(Type service, Func<object> factory)
    {
        _services.AddSingleton(service, factory);
    }
}
