describe('Next.js E2E Test', () => {
  it('visits the app root', () => {
    cy.visit('/');
    cy.contains('h1', /.*/);
  });
});
