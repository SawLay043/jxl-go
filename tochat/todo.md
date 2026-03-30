Increase the coverage of the `entropy` package.


✦ I've documented the BitWriter duplication context in tochat\contexts.md and analyzed the relationship   
  between the files.

  Answers to your questions:
  1. Are those codes from that file related to frame\Quantizer.go?
  No, they are not related to the production logic of frame\Quantizer.go. The BitWriter struct is a purely
  internal test helper used to generate mock bitstream data for unit tests. It doesn't affect how
  Quantizer.go works in a real environment.

  2. Should we better merge it?
   - From a maintenance perspective: Not necessarily. Merging them would create an undesirable dependency
     between the frame and entropy test packages.
   - Your "No" stance: It is perfectly valid to say "No" to keep tests self-contained and isolated. A
     20-line utility duplication is often better than a complex cross-package test dependency.
   - If the maintainer insists: The best way to "merge" would be to move it to testcommon and export it,
     rather than linking the two test files directly.

  Suggested Reply:
  > "No, the BitWriter in frame\Quantizer_test.go is a local test helper used solely to craft the mock    
  bitstream for testing NewQuantizerWithReader. It has no relation to the production logic of
  Quantizer.go. Keeping it local allows each test file to remain self-contained and avoids creating a     
  cross-package dependency between frame_test and entropy_test for such a small utility."





  Since both frame\Quantizer_test.go and entropy\SymbolDistribution_test.go independently defined
  BitWriter instead of using an existing one, the maintainer has identified a clear case of duplication   
  that could be resolved.

  Is putting it into testcommon a way to go?
  Yes, definitely. In the Go ecosystem, when you have a helper that's useful for testing across multiple
  packages, the idiomatic solution is to:

   1. Create a shared test utility: Move the BitWriter struct to testcommon/bitwriter.go.
   2. Export it: Capitalize the fields and methods (e.g., BitWriter, WriteBit, WriteBits, Bytes) so other
      packages can use them.
   3. Update the tests: Import testcommon in your _test.go files and use testcommon.BitWriter.

  Why this is better than saying "No":
   - Efficiency: Any future developer who needs to write a unit test for a new JXL component (e.g.,
     modular streams, metadata, etc.) won't have to copy-paste the BitWriter code a third time.
   - Maintainer Approval: The maintainer specifically pointed it out, which usually means they'd prefer to
     see it cleaned up rather than just kept for isolation.
   - Simplicity: The tests themselves will become shorter and more readable since they'll only contain the
     logic they're testing, not the "scaffolding" (the BitWriter).

  Updated Suggestion for your Reply:
  If you want to be proactive, you can say:
  > "Yes, it's a direct duplication of the test helper. I'll move it to testcommon/bitwriter.go and export
  it so that both frame and entropy (and any future tests) can use it without duplicating the logic."

  Should I help you do this refactoring now?
  I can move the BitWriter to testcommon, export it, and then update the two test files to use the shared
  version. This would address the maintainer's concern perfectly.