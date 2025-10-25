using Microsoft.Extensions.DependencyInjection;
using Spectre.Console.Cli;

namespace FolderManager;

public sealed class ServiceCollectionRegistrar(IServiceCollection services) : ITypeRegistrar
{
    public ITypeResolver Build()
    {
        return new ServiceCollectionResolver(services.BuildServiceProvider());
    }

    public void Register(Type service, Type implementation)
    {
        services.AddSingleton(service, implementation);
    }

    public void RegisterInstance(Type service, object implementation)
    {
        services.AddSingleton(service, implementation);
    }

    public void RegisterLazy(Type service, Func<object> factory)
    {
        services.AddSingleton(service, factory);
    }
}
