import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export var errorMessage = "amountExceedsBalance";
export function amountLessThanOrEqualToBalance(balance: number): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const amount = control.value;
    if (amount !== null && amount !== undefined && amount > balance) {
      return { amountExceedsBalance: true };
    }
    return null;
  };
}
