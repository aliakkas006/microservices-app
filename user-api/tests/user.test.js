const dummyUsers = [
  { id: 1, name: 'Ali Akkas', email: 'ali@example.com' },
  { id: 2, name: 'Rahim', email: 'rahim@example.com' },
  { id: 3, name: 'Karim', email: 'karim@example.com' },
];

describe('Dummy User Data', () => {
  test('should have 3 users', () => {
    expect(dummyUsers.length).toBe(3);
  });

  test('first user should be Ali Akkas', () => {
    expect(dummyUsers[0].name).toBe('Ali Akkas');
    expect(dummyUsers[0].email).toBe('ali@example.com');
  });

  test('second user should have valid email', () => {
    expect(dummyUsers[1].email).toMatch(/@/);
  });
});
