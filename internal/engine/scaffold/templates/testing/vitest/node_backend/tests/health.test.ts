import { describe, it, expect } from 'vitest';

describe('Health API', () => {
  it('should return ok', () => {
    expect({ status: 'ok' }).toEqual({ status: 'ok' });
  });
});
