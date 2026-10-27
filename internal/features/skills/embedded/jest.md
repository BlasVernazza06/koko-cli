---
name: jest
description: Jest test runner conventions, TypeScript setup (ts-jest/swc), isolated unit tests, mock restores, snapshot testing guidelines, and coverage metrics.
globs: ["**/*.{test,spec}.{ts,tsx,js,jsx}", "jest.config.*"]
---

# Production Jest Testing Standards & Architecture

Guidelines for deterministic, high-speed unit and integration testing with Jest, emphasizing mock hygiene, compiler acceleration (SWC), and behavioral assertions.

## Test Structure & Arrange-Act-Assert (AAA) Reference

| Phase | Purpose | Best Practice |
| :--- | :--- | :--- |
| **Arrange** | Setup inputs, mocks, test fixtures | Explicit local setup per test; avoid leaking state across tests |
| **Act** | Execute the function or method under test | Single line invocation when possible |
| **Assert** | Verify outputs, return values, errors | Assert expected behaviors rather than internal implementation details |

## Mock Hygiene Rules

- **Clear vs Reset vs Restore**:
  - `jest.clearAllMocks()`: Clears calls and results on mocks.
  - `jest.resetAllMocks()`: Clears calls, results, and resets mock implementations to `undefined`.
  - `jest.restoreAllMocks()`: Restores original non-mocked implementations (crucial for `jest.spyOn`).
- Always run `afterEach(() => { jest.restoreAllMocks(); });` to guarantee clean state.

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Fragile Snapshot Overuse & Unrestored Global Mocks
Relying on massive multi-page HTML snapshots creates noisy pull requests that developers blindly approve with `-u`.

```typescript
// BAD: Massive UI snapshot and uncontrolled global monkey-patching
test('renders entire user dashboard', () => {
  global.fetch = jest.fn().mockResolvedValue({ json: () => ({ id: '1' }) });
  const view = render(<FullAppDashboard />);
  expect(view.container).toMatchSnapshot(); // Huge 5,000-line snapshot that breaks on any CSS tweak!
});
```

### ✅ Recommended: Focused Behavioral Assertions & Safe Spies

```typescript
import { UserService } from './user.service';
import { UserRepository } from './user.repository';

describe('UserService', () => {
  let userService: UserService;
  let userRepoMock: jest.Mocked<UserRepository>;

  beforeEach(() => {
    userRepoMock = {
      findById: jest.fn(),
      save: jest.fn(),
    } as unknown as jest.Mocked<UserRepository>;

    userService = new UserService(userRepoMock);
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('should return user profile when user exists in repository', async () => {
    // Arrange
    const mockUser = { id: 'u-123', email: 'test@example.com', name: 'Alice' };
    userRepoMock.findById.mockResolvedValue(mockUser);

    // Act
    const result = await userService.getUser('u-123');

    // Assert
    expect(userRepoMock.findById).toHaveBeenCalledWith('u-123');
    expect(result).toEqual(mockUser);
  });

  it('should throw NotFoundError when user does not exist', async () => {
    // Arrange
    userRepoMock.findById.mockResolvedValue(null);

    // Act & Assert
    await expect(userService.getUser('non-existent')).rejects.toThrow('User not found');
  });
});
```

## Performance: Speeding up Jest with `@swc/jest`

Replace heavy `ts-jest` compilation with high-speed Rust-based SWC in `jest.config.ts`:

```typescript
// jest.config.ts
import type { Config } from 'jest';

const config: Config = {
  testEnvironment: 'node',
  transform: {
    '^.+\\.(t|j)sx?$': '@swc/jest',
  },
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 85,
      lines: 85,
      statements: 85,
    },
  },
};

export default config;
```
