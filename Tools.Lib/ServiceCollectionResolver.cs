using Spectre.Console.Cli;

namespace Tools.Lib;

public sealed class ServiceCollectionResolver(IServiceProvider provider) : ITypeResolver, IDisposable
{
    private readonly IServiceProvider _provider = provider ?? throw new ArgumentNullException(nameof(provider));

    public void Dispose()
    {
        if (_provider is IDisposable disposable)
        {
            disposable.Dispose();
        }
    }

    public object? Resolve(Type? type) => type is null ? null : _provider.GetService(type);
}
